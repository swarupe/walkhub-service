package walkhub

import (
	"errors"
	"net/http"

	"github.com/lib/pq"
	"github.com/tamasd/ab"
)

// AUTOGENERATED DO NOT EDIT

func NewEmbedLog() *EmbedLog {
	e := &EmbedLog{}

	// HOOK: newEmbedLog()

	return e
}

func EmptyEmbedLog() *EmbedLog {
	return &EmbedLog{}
}

var _ ab.Validator = &EmbedLog{}

func (e *EmbedLog) Validate() error {
	var err error

	err = validateEmbedLog(e)

	return err
}

func (e *EmbedLog) GetID() string {
	return e.UUID
}

var EmbedLogNotFoundError = errors.New("embedlog not found")

const embedlogFields = "e.uuid, e.ipaddr, e.created, e.site, e.mail"

func selectEmbedLogFromQuery(db ab.DB, query string, args ...interface{}) ([]*EmbedLog, error) {
	// HOOK: beforeEmbedLogSelect()

	entities := []*EmbedLog{}

	rows, err := db.Query(query, args...)

	if err != nil {
		return entities, err
	}

	for rows.Next() {
		e := EmptyEmbedLog()

		if err = rows.Scan(&e.UUID, &e.IPAddr, &e.Created, &e.Site, &e.Mail); err != nil {
			return []*EmbedLog{}, err
		}

		entities = append(entities, e)
	}

	// HOOK: afterEmbedLogSelect()

	return entities, err
}

func selectSingleEmbedLogFromQuery(db ab.DB, query string, args ...interface{}) (*EmbedLog, error) {
	entities, err := selectEmbedLogFromQuery(db, query, args...)
	if err != nil {
		return nil, err
	}

	if len(entities) > 0 {
		return entities[0], nil
	}

	return nil, nil
}

func (e *EmbedLog) Insert(db ab.DB) error {
	// HOOK: beforeEmbedLogInsert()

	err := db.QueryRow("INSERT INTO \"embedlog\"(ipaddr, created, site, mail) VALUES($1, $2, $3, $4) RETURNING uuid", e.IPAddr, e.Created, e.Site, e.Mail).Scan(&e.UUID)

	// HOOK: afterEmbedLogInsert()

	return err
}

func LoadEmbedLog(db ab.DB, UUID string) (*EmbedLog, error) {
	// HOOK: beforeEmbedLogLoad()

	e, err := selectSingleEmbedLogFromQuery(db, "SELECT "+embedlogFields+" FROM \"embedlog\" e WHERE e.uuid = $1", UUID)

	// HOOK: afterEmbedLogLoad()

	return e, err
}

func LoadAllEmbedLog(db ab.DB, start, limit int) ([]*EmbedLog, error) {
	// HOOK: beforeEmbedLogLoadAll()

	entities, err := selectEmbedLogFromQuery(db, "SELECT "+embedlogFields+" FROM \"embedlog\" e ORDER BY UUID DESC LIMIT $1 OFFSET $2", limit, start)

	// HOOK: afterEmbedLogLoadAll()

	return entities, err
}

type EmbedLogService struct {
}

func (s *EmbedLogService) Register(srv *ab.Server) error {
	var err error

	postMiddlewares := []func(http.Handler) http.Handler{}

	// HOOK: beforeEmbedLogServiceRegister()

	if err != nil {
		return err
	}

	srv.Post("/api/embedlog", s.embedlogPostHandler(), postMiddlewares...)

	// HOOK: afterEmbedLogServiceRegister()

	return err
}

func embedlogDBErrorConverter(err *pq.Error) ab.VerboseError {
	ve := ab.NewVerboseError(err.Message, err.Detail)

	// HOOK: convertEmbedLogDBError()

	return ve
}

func (s *EmbedLogService) embedlogPostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entity := &EmbedLog{}
		ab.MustDecode(r, entity)

		abort := false

		embedlogPostValidation(r, entity)

		if abort {
			return
		}

		if err := entity.Validate(); err != nil {
			ab.Fail(r, http.StatusBadRequest, err)
		}

		db := ab.GetDB(r)

		err := entity.Insert(db)
		ab.MaybeFail(r, http.StatusInternalServerError, ab.ConvertDBError(err, embedlogDBErrorConverter))

		afterEmbedLogPostInsertHandler(db, entity)

		if abort {
			return
		}

		ab.Render(r).SetCode(http.StatusCreated).JSON(entity)
	})
}

func (s *EmbedLogService) SchemaInstalled(db ab.DB) bool {
	found := ab.TableExists(db, "embedlog")

	// HOOK: afterEmbedLogSchemaInstalled()

	return found
}

func (s *EmbedLogService) SchemaSQL() string {
	sql := "CREATE TABLE \"embedlog\" (\n" +
		"\t\"uuid\" uuid DEFAULT uuid_generate_v4() NOT NULL,\n" +
		"\t\"ipaddr\" character varying NOT NULL,\n" +
		"\t\"created\" timestamp with time zone NOT NULL,\n" +
		"\t\"site\" character varying NOT NULL,\n" +
		"\t\"mail\" character varying NOT NULL,\n" +
		"\tCONSTRAINT embedlog_pkey PRIMARY KEY (uuid)\n);\n"

	// HOOK: afterEmbedLogSchemaSQL()

	return sql
}
