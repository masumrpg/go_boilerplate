package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type Config struct {
	Name        string // product
	NameUpper   string // Product
	NamePlural  string // products
	PackagePath string // go_boilerplate/internal/modules/product/dto
}

const (
	modulePath = "internal/modules"
	mainGoPath = "cmd/api/main.go"
)

var templates = map[string]string{
	"model.go": `package {{.Name}}

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// {{.NameUpper}} represents the {{.Name}} entity
type {{.NameUpper}} struct {
	ID        uuid.UUID      ` + "`" + `gorm:"type:uuid;primaryKey" json:"id"` + "`" + `
	Name      string         ` + "`" + `gorm:"type:varchar(255);not null" json:"name"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

// TableName overrides the table name used by {{.NameUpper}} to add prefix
func ({{.NameUpper}}) TableName() string {
	return "t_{{.NamePlural}}"
}
`,
	"repository.go": `package {{.Name}}

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type {{.NameUpper}}Repository interface {
	Create(item *{{.NameUpper}}) error
	FindByID(id uuid.UUID) (*{{.NameUpper}}, error)
	FindAll(page, limit int) ([]{{.NameUpper}}, int64, error)
	Update(item *{{.NameUpper}}) error
	Delete(id uuid.UUID) error
}

type {{.Name}}Repository struct {
	db *gorm.DB
}

func New{{.NameUpper}}Repository(db *gorm.DB) {{.NameUpper}}Repository {
	return &{{.Name}}Repository{db: db}
}

func (r *{{.Name}}Repository) Create(item *{{.NameUpper}}) error {
	return r.db.Create(item).Error
}

func (r *{{.Name}}Repository) FindByID(id uuid.UUID) (*{{.NameUpper}}, error) {
	var item {{.NameUpper}}
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *{{.Name}}Repository) FindAll(page, limit int) ([]{{.NameUpper}}, int64, error) {
	var items []{{.NameUpper}}
	var total int64
	offset := (page - 1) * limit

	if err := r.db.Model(&{{.NameUpper}}{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *{{.Name}}Repository) Update(item *{{.NameUpper}}) error {
	return r.db.Save(item).Error
}

func (r *{{.Name}}Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&{{.NameUpper}}{}, "id = ?", id).Error
}
`,
	"service.go": `package {{.Name}}

import (
	"{{.PackagePath}}"
	"github.com/google/uuid"
)

type {{.NameUpper}}Service interface {
	Create(req *dto.Create{{.NameUpper}}Request) (*{{.NameUpper}}, error)
	GetByID(id uuid.UUID) (*{{.NameUpper}}, error)
	GetAll(page, limit int) ([]{{.NameUpper}}, int64, error)
	Update(id uuid.UUID, req *dto.Update{{.NameUpper}}Request) (*{{.NameUpper}}, error)
	Delete(id uuid.UUID) error
}

type {{.Name}}Service struct {
	repo {{.NameUpper}}Repository
}

func New{{.NameUpper}}Service(repo {{.NameUpper}}Repository) {{.NameUpper}}Service {
	return &{{.Name}}Service{repo: repo}
}

func (s *{{.Name}}Service) Create(req *dto.Create{{.NameUpper}}Request) (*{{.NameUpper}}, error) {
	item := &{{.NameUpper}}{
		ID:   uuid.New(),
		Name: req.Name,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *{{.Name}}Service) GetByID(id uuid.UUID) (*{{.NameUpper}}, error) {
	return s.repo.FindByID(id)
}

func (s *{{.Name}}Service) GetAll(page, limit int) ([]{{.NameUpper}}, int64, error) {
	return s.repo.FindAll(page, limit)
}

func (s *{{.Name}}Service) Update(id uuid.UUID, req *dto.Update{{.NameUpper}}Request) (*{{.NameUpper}}, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		item.Name = req.Name
	}

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *{{.Name}}Service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
`,
	"handler.go": `package {{.Name}}

import (
	"strconv"

	"{{.PackagePath}}"
	"go_boilerplate/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type {{.NameUpper}}Handler struct {
	service {{.NameUpper}}Service
}

func New{{.NameUpper}}Handler(service {{.NameUpper}}Service) *{{.NameUpper}}Handler {
	return &{{.NameUpper}}Handler{service: service}
}

// Create handles creating a new {{.Name}}
// @Summary Create {{.Name}}
// @Tags {{.NameUpper}}
// @Accept json
// @Produce json
// @Param request body dto.Create{{.NameUpper}}Request true "Create data"
// @Success 201 {object} utils.APIResponse
// @Router /{{.NamePlural}} [post]
func (h *{{.NameUpper}}Handler) Create(c *fiber.Ctx) error {
	req := c.Locals("validatedBody").(*dto.Create{{.NameUpper}}Request)

	item, err := h.service.Create(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create {{.Name}}", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, item, "{{.NameUpper}} created successfully")
}

// Get handles retrieving a {{.Name}} by ID
// @Summary Get {{.Name}}
// @Tags {{.NameUpper}}
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} utils.APIResponse
// @Router /{{.NamePlural}}/{id} [get]
func (h *{{.NameUpper}}Handler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "{{.NameUpper}} not found", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, item, "{{.NameUpper}} retrieved successfully")
}

// List handles listing all {{.NamePlural}}
// @Summary List {{.NamePlural}}
// @Tags {{.NameUpper}}
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} utils.APIResponse
// @Router /{{.NamePlural}} [get]
func (h *{{.NameUpper}}Handler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	items, total, err := h.service.GetAll(page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve {{.NamePlural}}", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	}, "{{.NamePlural}} retrieved successfully")
}

// Update handles updating a {{.Name}}
// @Summary Update {{.Name}}
// @Tags {{.NameUpper}}
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param request body dto.Update{{.NameUpper}}Request true "Update data"
// @Success 200 {object} utils.APIResponse
// @Router /{{.NamePlural}}/{id} [put]
func (h *{{.NameUpper}}Handler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	req := c.Locals("validatedBody").(*dto.Update{{.NameUpper}}Request)

	item, err := h.service.Update(id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update {{.Name}}", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, item, "{{.NameUpper}} updated successfully")
}

// Delete handles deleting a {{.Name}}
// @Summary Delete {{.Name}}
// @Tags {{.NameUpper}}
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} utils.APIResponse
// @Router /{{.NamePlural}}/{id} [delete]
func (h *{{.NameUpper}}Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	if err := h.service.Delete(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete {{.Name}}", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "{{.NameUpper}} deleted successfully")
}
`,
	"routes.go": `package {{.Name}}

import (
	"{{.PackagePath}}"
	"go_boilerplate/internal/shared/config"
	"go_boilerplate/internal/shared/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, logger *logrus.Logger) {
	repo := New{{.NameUpper}}Repository(db)
	service := New{{.NameUpper}}Service(repo)
	handler := New{{.NameUpper}}Handler(service)

	api := app.Group("/api/v1/{{.NamePlural}}")
	api.Use(middleware.JWTAuth(cfg))

	api.Post("/", middleware.BodyValidator(&dto.Create{{.NameUpper}}Request{}), handler.Create)
	api.Get("/", handler.List)
	api.Get("/:id", handler.Get)
	api.Put("/:id", middleware.BodyValidator(&dto.Update{{.NameUpper}}Request{}), handler.Update)
	api.Delete("/:id", handler.Delete)
}
`,
	"dto/request.go": `package dto

type Create{{.NameUpper}}Request struct {
	Name string ` + "`" + `json:"name" validate:"required,min=3"` + "`" + `
}

type Update{{.NameUpper}}Request struct {
	Name string ` + "`" + `json:"name" validate:"omitempty,min=3"` + "`" + `
}
`,
	"dto/response.go": `package dto

import (
	"time"
	"github.com/google/uuid"
)

type {{.NameUpper}}Response struct {
	ID        uuid.UUID ` + "`" + `json:"id"` + "`" + `
	Name      string    ` + "`" + `json:"name"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}
`,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/gen/main.go <module-name>")
		os.Exit(1)
	}

	name := strings.ToLower(os.Args[1])
	nameUpper := strings.Title(name)
	namePlural := name + "s"
	if strings.HasSuffix(name, "y") {
		namePlural = name[:len(name)-1] + "ies"
	}

	config := Config{
		Name:        name,
		NameUpper:   nameUpper,
		NamePlural:  namePlural,
		PackagePath: "go_boilerplate/internal/modules/" + name + "/dto",
	}

	// 1. Create Directories
	baseDir := filepath.Join(modulePath, name)
	if err := os.MkdirAll(filepath.Join(baseDir, "dto"), 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// 2. Generate Files
	for fileName, tmplStr := range templates {
		filePath := filepath.Join(baseDir, fileName)

		tmpl, err := template.New(fileName).Parse(tmplStr)
		if err != nil {
			fmt.Printf("Error parsing template %s: %v\n", fileName, err)
			continue
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, config); err != nil {
			fmt.Printf("Error executing template %s: %v\n", fileName, err)
			continue
		}

		if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", filePath, err)
			continue
		}
		fmt.Printf("✓ Created %s\n", filePath)
	}

	// 3. Auto Inject to main.go
	injectToMain(config)

	// 4. Generate SQL Migrations
	generateMigrations(config)

	fmt.Printf("\n🚀 Module '%s' generated successfully!\n", name)
	fmt.Println("Next steps:")
	fmt.Printf("1. Refresh Swagger: make swagger\n")
}

func injectToMain(config Config) {
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		fmt.Printf("Error reading main.go: %v\n", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	for _, line := range lines {
		newLines = append(newLines, line)

		// Inject Import
		if strings.Contains(line, "// [MODULE_IMPORT_MARKER]") {
			newLines = append(newLines, fmt.Sprintf("\t%sModule \"go_boilerplate/internal/modules/%s\"", config.Name, config.Name))
		}

		// Inject Migration
		if strings.Contains(line, "// [MODULE_MIGRATION_MARKER]") {
			newLines = append(newLines, fmt.Sprintf("\t\t\t&%sModule.%s{},", config.Name, config.NameUpper))
		}

		// Inject Route
		if strings.Contains(line, "// [MODULE_ROUTE_MARKER]") {
			newLines = append(newLines, fmt.Sprintf("\t// %s routes", config.NameUpper))
			newLines = append(newLines, fmt.Sprintf("\t%sModule.RegisterRoutes(app, db, cfg, logger)", config.Name))
			newLines = append(newLines, fmt.Sprintf("\tlogger.Info(\"✓ %s routes registered\")", config.NameUpper))
			newLines = append(newLines, "")
		}
	}

	if err := os.WriteFile(mainGoPath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		fmt.Printf("Error updating main.go: %v\n", err)
	} else {
		fmt.Println("✓ Auto-injected to cmd/api/main.go")
	}
}

func generateMigrations(config Config) {
	timestamp := time.Now().Format("20060102150405")
	migrationDir := "db/migrations"

	upFileName := fmt.Sprintf("%s_create_%s_table.up.sql", timestamp, config.NamePlural)
	downFileName := fmt.Sprintf("%s_create_%s_table.down.sql", timestamp, config.NamePlural)

	upContent := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "t_%s" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS "idx_t_%s_deleted_at" ON "t_%s" ("deleted_at");
`, config.NamePlural, config.NamePlural, config.NamePlural)

	downContent := fmt.Sprintf(`DROP TABLE IF EXISTS "t_%s";
`, config.NamePlural)

	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		fmt.Printf("Error creating migration directory: %v\n", err)
		return
	}

	if err := os.WriteFile(filepath.Join(migrationDir, upFileName), []byte(upContent), 0644); err != nil {
		fmt.Printf("Error writing up migration: %v\n", err)
	} else {
		fmt.Printf("✓ Created %s/%s\n", migrationDir, upFileName)
	}

	if err := os.WriteFile(filepath.Join(migrationDir, downFileName), []byte(downContent), 0644); err != nil {
		fmt.Printf("Error writing down migration: %v\n", err)
	} else {
		fmt.Printf("✓ Created %s/%s\n", migrationDir, downFileName)
	}
}
