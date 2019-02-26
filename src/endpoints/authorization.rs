
use rocket_contrib::json::Json;

use rocket::{
	State, Response,
	http::Status
};

use crate::{
	DbConn, endpoints::{generate_error, generate_response},
	database::handler::HandlerError
};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct PostServiceAuth {
	pub refresh: Option<String>,
	pub expiration: Option<String>,
	pub access: String
}

#[get("/<channel>/<service>")]
pub fn get_service_auth<'r>(handler: State<DbConn>, channel: String, service: String) -> Response<'r> {
	let auth = handler.lock().expect("db lock").get_service_auth(&channel, &service);

	match auth {
		Ok(auth) => generate_response(Status::Ok, json!({
			"refresh": auth.refresh,
			"expiration": auth.expiration,
			"access": auth.access,
			"meta": json!({
				"service": service,
				"channel": channel				
			})
		})),
		Err(HandlerError::Error(_)) => generate_response(Status::NotFound, generate_error(404, None)),
		Err(e) => {
			println!("Internal error getting service auth: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[patch("/<channel>/<service>/update", format = "json", data = "<data>")]
pub fn update_service_auth<'r>(handler: State<DbConn>, channel: String, service: String, data: Json<PostServiceAuth>) -> Response<'r> {
	let result = handler.lock().expect("db lock").update_service_auth(&channel, &service, data.into_inner());

	match result {
		Ok(()) => generate_response(Status::Ok, json!({
			"updated": true,
			"meta": json!({
				"service": service,
				"channel": channel				
			})
		})),
		Err(HandlerError::Error(_)) => generate_response(Status::NotFound, generate_error(404, None)),
		Err(e) => {
			println!("Internal error getting service auth: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}
