
pub mod channel;

use rocket_contrib::json::JsonValue;

pub fn generate_error(code: u32) -> JsonValue {
	json!({
		"error": json! ({
			"code": code
		})
	})
}
