
use rocket_contrib::json::{JsonValue, Json};

use rocket::State;
use crate::{
	DbConn, endpoints::generate_error,
	database::{
		structures::Message,
		handler::HandlerError
	}
};

#[get("/<channel>/state")]
pub fn get_channel_state(handler: State<DbConn>, channel: String) -> JsonValue {
	let result = handler.lock().expect("db lock").get_channel_state(&channel, None);
	match result {
		Ok(state) => json! ({
			// TODO: Add more information
			"services": state.services,
			"meta": json!({
				"channel": state.token
			})
		}),
		Err(HandlerError::Error(e)) => generate_error(404, Some(e)),
		Err(e) => {
			println!("Internal error getting state: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[get("/<channel>/state/<service>")]
pub fn get_channel_service_state(handler: State<DbConn>, channel: String, service: String) -> JsonValue {
	let result = handler.lock().expect("db lock").get_channel_state(&channel, Some(service));
	match result {
		Ok(state) => json! ({
			// TODO: Add more information
			"service": state.services[0],
			"meta": json!({
				"channel": state.token
			})
		}),
		Err(HandlerError::Error(e)) => generate_error(404, Some(e)),
		Err(e) => {
			println!("Internal error getting state: {:?}", e);
			generate_error(500, None)
		}
	}
}
