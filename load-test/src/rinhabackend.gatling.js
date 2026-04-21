import {
  simulation,
  scenario,
  exec,
  csv,
  pause,
  css,
  feed,
  repeat,
  tryMax,
  rampUsers,
  tsv,
  StringBody,
  constantUsersPerSec,
  rampUsersPerSec,
} from "@gatling.io/core";
import { http, status, header } from "@gatling.io/http";

export default simulation((setUp) => {
  const http_protocol = http
    .baseUrl("http://localhost:9999")
    .userAgentHeader("xama na cockfight");

  const cria_consulta = scenario("criação e talvez consulta de pessoas")
    .feed(tsv("pessoas-payloads.tsv").circular())
    .exec(
      http("criação")
        .post("/pessoas")
        .body(StringBody("#{payload}"))
        .header("Content-Type", "application/json")
        .check(status().in(201, 422, 400))
        .check(status().saveAs("httpStatus"))
        .checkIf((session) => session.get("httpStatus") == "201")
        .then(header("Location").saveAs("location")),
    )
    .pause(
      { amount: 1, unit: "milliseconds" },
      { amount: 30, unit: "milliseconds" },
    )
    .doIf((session) => session.contains("location"))
    .then(http("consulta").get("/pessoas/#{location}"));

  const busca = scenario("busca válida de pessoas")
    .feed(tsv("termos-busca.tsv").circular())
    .exec(http("busca válida").get("/pessoas?t=#{t}"));

  const busca_invalida = scenario("busca inválida de pessoas").exec(
    http("busca inválida").get("/pessoas").check(status().is(400)),
  );

  const rampCria = 600;
  const rampBusca = 100;
  const rampBuscaInvalida = 40;

  // const rampCria = 1200
  // const rampBusca = 200
  // const rampBuscaInvalida = 80

  setUp(
    cria_consulta.injectOpen(
      constantUsersPerSec(2).during({ amount: 10, unit: "seconds" }),
      constantUsersPerSec(5).during(15).randomized(),

      rampUsersPerSec(6).to(rampCria).during({ amount: 3, unit: "minutes" }),
    ),
    busca.injectOpen(
      constantUsersPerSec(2).during({ amount: 25, unit: "seconds" }),

      rampUsersPerSec(6).to(100).during({ amount: 3, unit: "minutes" }),
    ),
    busca_invalida.injectOpen(
      constantUsersPerSec(2).during({ amount: 25, unit: "seconds" }),

      // rampUsersPerSec(6).to(40).during({ amount: 3, unit: 'minutes' }),
      rampUsersPerSec(6).to(80).during({ amount: 3, unit: "minutes" }),
    ),
  ).protocols(http_protocol);
});
