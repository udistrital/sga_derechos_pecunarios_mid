package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_derecho_pecunario_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

type DerechosPecuniariosController struct {
	beego.Controller
}

func (c *DerechosPecuniariosController) URLMapping() {
	c.Mapping("PostConcepto", c.PostConcepto)
	c.Mapping("PutConcepto", c.PutConcepto)
	c.Mapping("PostClonarConceptos", c.PostClonarConceptos)
	c.Mapping("GetDerechosPecuniariosPorVigencia", c.GetDerechosPecuniariosPorVigencia)
	c.Mapping("DeleteConcepto", c.DeleteConcepto)
	c.Mapping("PutCostoConcepto", c.PutCostoConcepto)
	c.Mapping("PostGenerarDerechoPecuniarioEstudiante", c.PostGenerarDerechoPecuniarioEstudiante)
	c.Mapping("GetEstadoRecibo", c.GetEstadoRecibo)
	c.Mapping("GetConsultarPersona", c.GetConsultarPersona)
	c.Mapping("PostSolicitudDerechoPecuniario", c.PostSolicitudDerechoPecuniario)
	c.Mapping("GetSolicitudDerechoPecuniario", c.GetSolicitudDerechoPecuniario)
	c.Mapping("PostRespuestaSolicitudDerechoPecuniario", c.PostRespuestaSolicitudDerechoPecuniario)
}

// PostConcepto ...
// @Title PostConcepto
// @Description Agregar un concepto
// @Param	body		body 	{}	true		"body Agregar Concepto content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /conceptos [post]
func (c *DerechosPecuniariosController) PostConcepto() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	resultado, err := services.PostConcepto(data)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// PutConcepto ...
// @Title PutConcepto
// @Description Modificar un concepto
// @Param	body		body 	{}	true		"body Modificar Concepto content"
// @Success 200 {}
// @Failure 400 body is empty
// @Failure 404 no data found
// @Failure 403 :id is empty
// @router /conceptos/:id [put]
func (c *DerechosPecuniariosController) PutConcepto() {
	defer errorhandler.HandlePanic(&c.Controller)

	idConcepto := c.Ctx.Input.Param(":id")
	data := c.Ctx.Input.RequestBody

	resultado, err := services.PutConcepto(idConcepto, data)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// DeleteConcepto ...
// @Title DeleteConcepto
// @Description Inactivar Concepto y Factor por id
// @Param   id      path    string  true        "Id del Concepto"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /conceptos/:id [delete]
func (c *DerechosPecuniariosController) DeleteConcepto() {
	defer errorhandler.HandlePanic(&c.Controller)

	idConcepto := c.Ctx.Input.Param(":id")

	resultado, err := services.DeleteConcepto(idConcepto)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// GetDerechosPecuniariosPorVigencia ...
// @Title GetDerechosPecuniariosPorVigencia
// @Description Consulta los derechos pecuniarias de la vigencia por id
// @Param	id		path	int	true	"Id de la vigencia correspondiente"
// @Success 200 {}
// @Failure 403 :id is empty
// @Failure 404 no data found
// @router /vigencias/:id [get]
func (c *DerechosPecuniariosController) GetDerechosPecuniariosPorVigencia() {
	defer errorhandler.HandlePanic(&c.Controller)

	idVigencia := c.Ctx.Input.Param(":id")

	resultado, err := services.GetDerechosPecuniariosPorVigencia(idVigencia)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// PostClonarConceptos ...
// @Title PostClonarConceptos
// @Description Clona los conceptos de la vigencia anterior en la vigencia actual
// @Param	body		body 	{}	true		"body Clonar Conceptos content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /vigencias/clonar-conceptos [post]
func (c *DerechosPecuniariosController) PostClonarConceptos() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	resultado, err := services.PostClonarConceptos(data)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// PutCostoConcepto ...
// @Title PutCostoConcepto
// @Description Añadir el costo de un concepto existente
// @Param   body        body    {}  true        "body Inhabilitar Proyecto content"
// @Success 200 {}
// @Failure 400 :body is empty
// @router /conceptos/costos [post]
func (c *DerechosPecuniariosController) PutCostoConcepto() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	resultado, err := services.PutCostoConcepto(data)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// PostGenerarDerechoPecuniarioEstudiante ...
// @Title PostGenerarrDerechoPecuniarioEstudiante
// @Description Generar un recibo de derecho pecuniario por parte de estudiantes
// @Param	body		body 	{}	true		"body Clonar Conceptos content"
// @Success 200 {}
// @Failure 404 not found resource
// @Failure 400 body is empty
// @router /derechos [post]
func (c *DerechosPecuniariosController) PostGenerarDerechoPecuniarioEstudiante() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	resultado, err := services.PostGenerarDerechoPecuniarioEstudiante(data)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// GetEstadoRecibo ...
// @Title GetEstadoRecibo
// @Description consultar los estados de todos los recibos de derechos pecuniarios generados por el tercero
// @Param	persona_id	path	int	true	"Id del tercero"
// @Param	periodo_id	path	int	true	"Id del ultimo periodo"
// @Success 200 {}
// @Failure 404 not found resource
// @router /personas/:persona_id/periodos/:periodo_id/recibos [get]
func (c *DerechosPecuniariosController) GetEstadoRecibo() {
	defer errorhandler.HandlePanic(&c.Controller)

	idPersona := c.Ctx.Input.Param(":persona_id")
	idPeriodo := c.Ctx.Input.Param(":periodo_id")

	resultado, err := services.GetEstadoRecibo(idPersona, idPeriodo)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// GetConsultarPersona ...
// @Title GetConsultarPersona
// @Description get información del estudainte por el id de tercero
// @Param	persona_id	path	int	true	"Id del tercero"
// @Success 200 {}
// @Failure 404 not found resource
// @router /personas/:persona_id [get]
func (c *DerechosPecuniariosController) GetConsultarPersona() {
	defer errorhandler.HandlePanic(&c.Controller)

	//Id del tercero
	idPersona := c.Ctx.Input.Param(":persona_id")

	resultado, err := services.GetConsultarPersona(idPersona)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// PostSolicitudDerechoPecuniario ...
// @Title PostSolicitudDerechoPecuniario
// @Description Crear una solicitud de derecho pecuniario
// @Param	body		body 	{}	true		"body Agregar Concepto content"
// @Success 200 {}
// @Failure 400 body is empty
// @router /solicitudes [post]
func (c *DerechosPecuniariosController) PostSolicitudDerechoPecuniario() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	resultado, err := services.PostSolicitudDerechoPecuniario(data)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// GetSolicitudDerechoPecuniario ...
// @Title GetSolicitudDerechoPecuniario
// @Description Obtener todos las solicitudes de derechos pecuniarios
// @Success 200 {}
// @Failure 400 body is empty
// @router /solicitudes [get]
func (c *DerechosPecuniariosController) GetSolicitudDerechoPecuniario() {
	defer errorhandler.HandlePanic(&c.Controller)

	resultado, err := services.GetSolicitudDerechoPecuniario()

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// PostRespuestaSolicitudDerechoPecuniario ...
// @Title PostRespuestaSolicitudDerechoPecuniario
// @Description Da respuesta a la solicitud de derechos pecuniarios
// @Param	id	path	int	true	"Id de la solicitud"
// @Success 200 {}
// @Failure 400 body is empty
// @router /solicitudes/:id/respuesta [post]
func (c *DerechosPecuniariosController) PostRespuestaSolicitudDerechoPecuniario() {
	defer errorhandler.HandlePanic(&c.Controller)

	idSolicitud := c.Ctx.Input.Param(":id")

	data := c.Ctx.Input.RequestBody

	resultado, err := services.PostRespuestaSolicitudDerechoPecuniario(idSolicitud, data)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

	c.ServeJSON()
}
