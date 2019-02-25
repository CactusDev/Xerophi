
pub mod channel;
pub mod state;
pub mod authorization;
pub mod quote;

use rocket::{Response, http::{Status, ContentType}};
use rocket_contrib::json::JsonValue;
use std::io::Cursor;

pub fn generate_error(code: u32, message: Option<String>) -> JsonValue {
	json!({
		"error": json! ({
			"code": code,
			"message": message.unwrap_or("".to_string())
		})
	})
}

pub fn generate_response<'r>(code: Status, json: JsonValue) -> Response<'r> {
	let mut response = Response::build();
	let json = json.to_string();

    response.status(code);
	response.sized_body(Cursor::new(json));
	response.header(ContentType::JSON);

	response.finalize()
}
