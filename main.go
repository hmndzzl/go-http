package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Estructura de las bandas
type Banda struct {
	ID             int    `json:"id"`
	Nombre         string `json:"nombre"`
	Genero         string `json:"genero"`
	PaisOrigen     string `json:"pais_origen"`
	AnioFormacion  int    `json:"anio_formacion"`
	AlbumDestacado string `json:"album_destacado"`
}

// Estructura para respuestas de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Variable global para almacenar las bandas en memoria
var bandas []Banda

func main() {
	loadBandas() // Cargar bandas

	mux := http.NewServeMux() // mux para manejar rutas con parámetros

	mux.HandleFunc("GET /api/ping", pingHandler) // Endpoint de prueba
	mux.HandleFunc("GET /api/bandas", handleGetBandas) // Endpoint para hecer GET /api/bandas
	mux.HandleFunc("POST /api/bandas", handleCreateBanda) // Endpoint para hacer POST /api/bandas

	mux.HandleFunc("GET /api/bandas/{id}", handleGetBandaByID) // Endpoint para hacer GET /api/bandas/{id} con id
	mux.HandleFunc("PUT /api/bandas/{id}", handleUpdateBanda) // Endpoint para hacer PUT para actualizar una banda con id
	mux.HandleFunc("PATCH /api/bandas/{id}", handlePatchBanda) // Enpoint para hacer PATCH /api/bandas/{id}
	mux.HandleFunc("DELETE /api/bandas/{id}", handleDeleteBanda) // Enpoint para borrar una banda por id

	log.Println("POST JSON API running on :41265")
	log.Fatal(http.ListenAndServe(":41265", mux))
}

func loadBandas() {
	file, err := os.ReadFile("./data/bandas.json")
	if err != nil {
		log.Fatal("Error leyendo archivo:", err)
	}
	if err = json.Unmarshal(file, &bandas); err != nil {
		log.Fatal("Error parseando JSON:", err)
	}
}

func saveBandas() error {
	data, err := json.MarshalIndent(bandas, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("./data/bandas.json", data, 0644)
}

func generateNextID() int {
	maxID := 0
	for _, b := range bandas {
		if b.ID > maxID {
			maxID = b.ID
		}
	}
	return maxID + 1
}

type apiError struct{ msg string }

func (e *apiError) Error() string { return e.msg }

func validateBanda(b Banda) error {
	if strings.TrimSpace(b.Nombre) == "" {
		return &apiError{"El campo 'nombre' es requerido"}
	}
	if strings.TrimSpace(b.Genero) == "" {
		return &apiError{"El campo 'genero' es requerido"}
	}
	if strings.TrimSpace(b.PaisOrigen) == "" {
		return &apiError{"El campo 'pais_origen' es requerido"}
	}
	if b.AnioFormacion < 1900 || b.AnioFormacion > 2030 {
		return &apiError{"El campo 'anio_formacion' debe ser un año válido entre 1900 y 2030"}
	}
	if strings.TrimSpace(b.AlbumDestacado) == "" {
		return &apiError{"El campo 'album_destacado' es requerido"}
	}
	return nil
}

// Handlers

func pingHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "pong"})
}

func handleGetBandas(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if idParam := q.Get("id"); idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_PARAM", "El parámetro 'id' debe ser un número entero")
			return
		}

		for _, b := range bandas {
			if b.ID == id {
				writeJSON(w, http.StatusOK, b)
				return
			}
		}
		writeError(w, http.StatusNotFound, "NOT_FOUND", "No se encontró una banda con ese id")
		return
	}

	nombre := strings.ToLower(q.Get("nombre"))
	genero := strings.ToLower(q.Get("genero"))
	pais := strings.ToLower(q.Get("pais_origen"))
	anioDesdeStr := q.Get("anio_desde")
	anioHastaStr := q.Get("anio_hasta")

	var anioDesde, anioHasta int
	if anioDesdeStr != "" {
		v, err := strconv.Atoi(anioDesdeStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_PARAM", "El parámetro 'anio_desde' debe ser un número entero")
			return
		}
		anioDesde = v
	}
	if anioHastaStr != "" {
		v, err := strconv.Atoi(anioHastaStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_PARAM", "El parámetro 'anio_hasta' debe ser un número entero")
			return
		}
		anioHasta = v
	}


	result := []Banda{}
	for _, b := range bandas {
		if nombre != "" && !strings.Contains(strings.ToLower(b.Nombre), nombre) {
			continue
		}
		if genero != "" && !strings.Contains(strings.ToLower(b.Genero), genero) {
			continue
		}
		if pais != "" && !strings.Contains(strings.ToLower(b.PaisOrigen), pais) {
			continue
		}
		if anioDesde > 0 && b.AnioFormacion < anioDesde {
			continue
		}
		if anioHasta > 0 && b.AnioFormacion > anioHasta {
			continue
		}
		result = append(result, b)
	}

	writeJSON(w, http.StatusOK, result)
}

func handleGetBandaByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PARAM", "El id debe ser un número entero")
		return
	}

	for _, b := range bandas {
		if b.ID == id {
			writeJSON(w, http.StatusOK, b)
			return
		}
	}
	writeError(w, http.StatusNotFound, "NOT_FOUND", "No se encontró una banda con ese id")
}

func handleCreateBanda(w http.ResponseWriter, r *http.Request) {
	var nueva Banda
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "El body no contiene JSON válido")
		return
	}
	if err := validateBanda(nueva); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	nueva.ID = generateNextID()
	bandas = append(bandas, nueva)
	if err := saveBandas(); err != nil {
		log.Println("Error guardando bandas:", err)
	}
	writeJSON(w, http.StatusCreated, nueva)
}

func handleUpdateBanda(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PARAM", "El id debe ser un número entero")
		return
	}

	var updated Banda
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "El body no contiene JSON válido")
		return
	}
	if err := validateBanda(updated); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	for i, b := range bandas {
		if b.ID == id {
			updated.ID = id
			bandas[i] = updated
			if err := saveBandas(); err != nil {
				log.Println("Error guardando bandas:", err)
			}
			writeJSON(w, http.StatusOK, updated)
			return
		}
	}
	writeError(w, http.StatusNotFound, "NOT_FOUND", "No se encontró una banda con ese id")
}

func handlePatchBanda(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PARAM", "El id debe ser un número entero")
		return
	}

	var patch map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "El body no contiene JSON válido")
		return
	}
	if len(patch) == 0 {
		writeError(w, http.StatusBadRequest, "EMPTY_BODY", "El body no puede estar vacío")
		return
	}

	for i, b := range bandas {
		if b.ID == id {
			if v, ok := patch["nombre"]; ok {
				s, ok := v.(string)
				if !ok || strings.TrimSpace(s) == "" {
					writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "El campo 'nombre' no puede estar vacío")
					return
				}
				bandas[i].Nombre = s
			}
			if v, ok := patch["genero"]; ok {
				s, ok := v.(string)
				if !ok || strings.TrimSpace(s) == "" {
					writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "El campo 'genero' no puede estar vacío")
					return
				}
				bandas[i].Genero = s
			}
			if v, ok := patch["pais_origen"]; ok {
				s, ok := v.(string)
				if !ok || strings.TrimSpace(s) == "" {
					writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "El campo 'pais_origen' no puede estar vacío")
					return
				}
				bandas[i].PaisOrigen = s
			}
			if v, ok := patch["anio_formacion"]; ok {
				f, ok := v.(float64)
				if !ok {
					writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "El campo 'anio_formacion' debe ser un número")
					return
				}
				year := int(f)
				if year < 1900 || year > 2026 {
					writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "El campo 'anio_formacion' debe ser un año válido entre 1900 y 2026")
					return
				}
				bandas[i].AnioFormacion = year
			}
			if v, ok := patch["album_destacado"]; ok {
				s, ok := v.(string)
				if !ok || strings.TrimSpace(s) == "" {
					writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "El campo 'album_destacado' no puede estar vacío")
					return
				}
				bandas[i].AlbumDestacado = s
			}
			if err := saveBandas(); err != nil {
				log.Println("Error guardando bandas:", err)
			}
			writeJSON(w, http.StatusOK, bandas[i])
			return
		}
	}
	writeError(w, http.StatusNotFound, "NOT_FOUND", "No se encontró una banda con ese id")
}

func handleDeleteBanda(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PARAM", "El id debe ser un número entero")
		return
	}

	for i, b := range bandas {
		if b.ID == id {
			bandas = append(bandas[:i], bandas[i+1:]...)
			if err := saveBandas(); err != nil {
				log.Println("Error guardando bandas:", err)
			}
			writeJSON(w, http.StatusOK, map[string]string{"message": "Banda eliminada correctamente"})
			return
		}
	}
	writeError(w, http.StatusNotFound, "NOT_FOUND", "No se encontró una banda con ese id")
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{
		Error:   code,
		Code:    status,
		Message: message,
	})
}