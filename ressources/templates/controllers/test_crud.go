package controllers

import (
	"github.com/solher/zest/ressources/templates"
	"github.com/clipperhouse/typewriter"
)

func init() {
	imports := []typewriter.ImportSpec{
		typewriter.ImportSpec{
			Name: ".",
			Path: "github.com/smartystreets/goconvey/convey",
		},
	}

	err := typewriter.Register(templates.NewWrite("controller_test", testSlice, imports))
	if err != nil {
		panic(err)
	}
}

var testSlice = typewriter.TemplateSlice{
	controllerTest,
}

var controllerTest = &typewriter.Template{
	Name: "Controller_test",
	Text: `
	type stubRender struct {
		renderer *render.Render
	}

	func newStubRender() *stubRender {
		return &stubRender{renderer: render.New()}
	}

	func (r *stubRender) JSONError(w http.ResponseWriter, status int, apiError *apierrors.APIError, err error) {
		r.renderer.JSON(w, status, apierrors.Make(*apiError, status, err))
	}

	func (r *stubRender) JSON(w http.ResponseWriter, status int, object interface{}) {
		r.renderer.JSON(w, status, object)
	}

	type stubInteractor struct {
		ThrowErrors bool
	}

	func (i *stubInteractor) Create({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
		if i.ThrowErrors {
			return nil, errors.New("error")
		}

		for i := range {{.Name}}s {
			{{.Name}}s[i].ID = i + 1
		}

		return {{.Name}}s, nil
	}

	func (i *stubInteractor) Find(filter *interfaces.Filter) ([]{{.Type}}, error) {
		{{.Name}}s := []{{.Type}}{
			{
				FirstName: "Fabien",
				LastName:  "Herfray",
			},
			{
				FirstName: "Thomas",
				LastName:  "Hourlier",
			},
		}

		if i.ThrowErrors {
			return nil, errors.New("error")
		}

		return {{.Name}}s, nil
	}

	func (i *stubInteractor) FindByID(id int, filter *interfaces.Filter) (*{{.Type}}, error) {
		{{.Name}}s := []{{.Type}}{
			{
				FirstName: "Fabien",
				LastName:  "Herfray",
			},
		}

		if i.ThrowErrors || id-1 > 0 {
			return nil, errors.New("error")
		}

		return &{{.Name}}s[id-1], nil
	}

	func (i *stubInteractor) Upsert({{.Name}}s []{{.Type}}) ([]{{.Type}}, error) {
		if i.ThrowErrors {
			return nil, errors.New("error")
		}

		for i := range {{.Name}}s {
			if {{.Name}}s[i].ID == 0 {
				switch {
				case {{.Name}}s[i-1].ID != 0:
					{{.Name}}s[i].ID = {{.Name}}s[i-1].ID
				default:
					{{.Name}}s[i].ID = i + 1
				}
			}
		}

		return {{.Name}}s, nil
	}

	func (i *stubInteractor) DeleteAll(filter *interfaces.Filter) error {
		if i.ThrowErrors {
			return errors.New("error")
		}

		return nil
	}

	func (i *stubInteractor) DeleteByID(id int) error {
		if i.ThrowErrors {
			return errors.New("error")
		}

		return nil
	}

	func TestController(t *testing.T) {
		interactor := &stubInteractor{}
		render := infrastructure.NewRender()
		routes := interfaces.NewRouteDirectory()

		controller := NewController(interactor, render, routes)
		key := usecases.NewDirectoryKey(controller)

		Convey("Testing {{.Name}}s controller...", t, func() {
			{{.Name}}s := []{{.Type}}{
				{
					FirstName: "Fabien",
					LastName:  "Herfray",
				},
				{
					FirstName: "Thomas",
					LastName:  "Hourlier",
				},
			}

			interactor.ThrowErrors = false

			Convey("Should be able to create {{.Name}}s.", func() {
				route := routes[key.For("Create")]
				res := interfaces.MockHTTPRequest(route, utils.MarshalToStr({{.Name}}s), "", nil)
				utils.Unmarshal(res, &{{.Name}}s)

				So({{.Name}}s[0].ID, ShouldNotEqual, 0)
				So({{.Name}}s[1].ID, ShouldNotEqual, 0)

				res = interfaces.MockHTTPRequest(route, "string", "", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.BodyDecodingError.ErrorCode)

				res = interfaces.MockHTTPRequest(route, "", "", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.BodyDecodingError.ErrorCode)

				interactor.ThrowErrors = true

				res = interfaces.MockHTTPRequest(route, utils.MarshalToStr({{.Name}}s), "", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.InternalServerError.ErrorCode)
			})

			Convey("Should be able to find {{.Name}}s.", func() {
				route := routes[key.For("Find")]

				res := interfaces.MockHTTPRequest(route, "", "", nil)
				So(res, ShouldEqual, utils.MarshalToStr({{.Name}}s))

				res = interfaces.MockHTTPRequest(route, "", "toto", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.FilterDecodingError.ErrorCode)

				res = interfaces.MockHTTPRequest(route, "", "{\"limit\": 1}", nil)
				So(interfaces.GetErrorCode(res), ShouldNotEqual, apierrors.FilterDecodingError.ErrorCode)

				interactor.ThrowErrors = true

				res = interfaces.MockHTTPRequest(route, "", "", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.InternalServerError.ErrorCode)
			})

			Convey("Should be able to find {{.Name}} by id.", func() {
				route := routes[key.For("FindByID")]

				params := map[string]string{
					{
						Key:   "id",
						Value: "1",
					},
				}

				res := interfaces.MockHTTPRequest(route, "", "", params)
				So(res, ShouldEqual, utils.MarshalToStr({{.Name}}s[0]))

				params[0].Value = "2"
				res = interfaces.MockHTTPRequest(route, "", "", params)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.Unauthorized.ErrorCode)

				params[0].Value = ""
				res = interfaces.MockHTTPRequest(route, "", "", params)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.InvalidPathParams.ErrorCode)

				params[0].Value = "toto"
				res = interfaces.MockHTTPRequest(route, "", "", params)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.InvalidPathParams.ErrorCode)

				interactor.ThrowErrors = true

				params[0].Value = "1"
				res = interfaces.MockHTTPRequest(route, "", "", params)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.Unauthorized.ErrorCode)
			})

			Convey("Should be able to upsert {{.Name}}s.", func() {
				route := routes[key.For("Upsert")]
				{{.Name}}s[0].ID = 3

				res := interfaces.MockHTTPRequest(route, utils.MarshalToStr({{.Name}}s), "", nil)
				utils.Unmarshal(res, &{{.Name}}s)

				So({{.Name}}s[0].ID, ShouldEqual, 3)
				So({{.Name}}s[1].ID, ShouldNotEqual, 0)

				res = interfaces.MockHTTPRequest(route, "string", "", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.BodyDecodingError.ErrorCode)

				res = interfaces.MockHTTPRequest(route, "", "", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.BodyDecodingError.ErrorCode)

				interactor.ThrowErrors = true

				res = interfaces.MockHTTPRequest(route, utils.MarshalToStr({{.Name}}s), "", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.InternalServerError.ErrorCode)
			})

			Convey("Should be able to delete {{.Name}}s.", func() {
				route := routes[key.For("DeleteAll")]

				res := interfaces.MockHTTPRequest(route, utils.MarshalToStr({{.Name}}s), "", nil)
				So(res, ShouldEqual, "null")

				res = interfaces.MockHTTPRequest(route, "", "toto", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.FilterDecodingError.ErrorCode)

				res = interfaces.MockHTTPRequest(route, "", "{\"limit\": 1}", nil)
				So(interfaces.GetErrorCode(res), ShouldNotEqual, apierrors.FilterDecodingError.ErrorCode)

				interactor.ThrowErrors = true

				res = interfaces.MockHTTPRequest(route, utils.MarshalToStr({{.Name}}s), "", nil)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.InternalServerError.ErrorCode)
			})

			Convey("Should be able to delete {{.Name}} by id.", func() {
				route := routes[key.For("DeleteByID")]

				params := map[string]string{
					{
						Key:   "id",
						Value: "1",
					},
				}

				res := interfaces.MockHTTPRequest(route, "", "", params)
				So(res, ShouldEqual, "null")

				params[0].Value = "2"
				res = interfaces.MockHTTPRequest(route, "", "", params)
				So(res, ShouldEqual, "null")

				params[0].Value = ""
				res = interfaces.MockHTTPRequest(route, "", "", params)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.InvalidPathParams.ErrorCode)

				params[0].Value = "toto"
				res = interfaces.MockHTTPRequest(route, "", "", params)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.InvalidPathParams.ErrorCode)

				interactor.ThrowErrors = true

				params[0].Value = "1"
				res = interfaces.MockHTTPRequest(route, "", "", params)
				So(interfaces.GetErrorCode(res), ShouldEqual, apierrors.Unauthorized.ErrorCode)
			})
		})
	}
`}
