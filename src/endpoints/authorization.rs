
use rocket_contrib::json::{JsonValue, Json};

use rocket::State;
use crate::{
	DbConn, endpoints::generate_error,
	database::handler::HandlerError
};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct PostServiceAuth {
	pub refresh: Option<String>,
	pub expiration: Option<String>,
	pub access: String
}

#[get("/<channel>/<service>")]
pub fn get_service_auth(handler: State<DbConn>, channel: String, service: String) -> JsonValue {
	let auth = handler.lock().expect("db lock").get_service_auth(&channel, &service);

	match auth {
		Ok(auth) => json!({
			"refresh": auth.refresh,
			"expiration": auth.expiration,
			"access": auth.access,
			"meta": json!({
				"service": service,
				"channel": channel				
			})
		}),
		Err(HandlerError::Error(e)) => generate_error(404, Some(e)),
		Err(e) => {
			println!("Internal error getting service auth: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[patch("/<channel>/<service>/update", format = "json", data = "<data>")]
pub fn update_service_auth(handler: State<DbConn>, channel: String, service: String, data: Json<PostServiceAuth>) -> JsonValue {
	let result = handler.lock().expect("db lock").update_service_auth(&channel, &service, data.into_inner());

	match result {
		Ok(()) => json!({
			"updated": true,
			"meta": json!({
				"service": service,
				"channel": channel				
			})
		}),
		Err(HandlerError::Error(e)) => generate_error(404, Some(e)),
		Err(e) => {
			println!("Internal error getting service auth: {:?}", e);
			generate_error(500, None)
		}
	}
}
