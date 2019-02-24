
pub mod channel;
pub mod state;

use rocket_contrib::json::JsonValue;

pub fn generate_error(code: u32, message: Option<String>) -> JsonValue {
	json!({
		"error": json! ({
			"code": code,
			"message": message.unwrap_or("".to_string())
		})
	})
}
