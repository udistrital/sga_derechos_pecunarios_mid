package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostConcepto",
            Router: "/conceptos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PutConcepto",
            Router: "/conceptos/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "DeleteConcepto",
            Router: "/conceptos/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PutCostoConcepto",
            Router: "/conceptos/costo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostGenerarDerechoPecuniarioEstudiante",
            Router: "/derechos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "GetConsultarPersona",
            Router: "/personas/:persona_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "GetEstadoRecibo",
            Router: "/personas/:persona_id/periodos/:id_periodo/estados_recibos",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostSolicitudDerechoPecuniario",
            Router: "/solicitudes",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "GetSolicitudDerechoPecuniario",
            Router: "/solicitudes",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostRespuestaSolicitudDerechoPecuniario",
            Router: "/solicitudes/:id/respuesta",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "GetDerechosPecuniariosPorVigencia",
            Router: "/vigencias/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_mid_derechos_pecunarios/controllers:DerechosPecuniariosController"],
        beego.ControllerComments{
            Method: "PostClonarConceptos",
            Router: "/vigencias/clonar_conceptos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
