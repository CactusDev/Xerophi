
use rocket_contrib::json::{JsonValue, Json};

use rocket::State;
use crate::{
	DbConn, endpoints::generate_error,
	database::{
		handler::HandlerError,
		structures::Message
	}
};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct PostQuote {
	pub response: Vec<Message>
}

#[get("/<channel>")]
pub fn get_quote(handler: State<DbConn>, channel: String) -> JsonValue {
	let quotes = handler.lock().expect("db lock").get_quote(&channel, None);
	match quotes {
		Ok(quotes) => json!({ "data": quotes }),
		Err(HandlerError::Error(_)) => json!({ "data": json!([]) }),
		Err(e) => {
			println!("Internal error getting quote: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[get("/<channel>/<id>", rank = 3)]
pub fn get_quote_by_id(handler: State<DbConn>, channel: String, id: u32) -> JsonValue {
	let quote = handler.lock().expect("db lock").get_quote(&channel, Some(id));
	match quote {
		Ok(quote) => json!({ "data": quote[0] }),
		Err(HandlerError::Error(_)) => json!({ "data": json!({}) }),
		Err(e) => {
			println!("Internal error getting quote: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[post("/<channel>/create", rank = 2, format = "json", data = "<quote>")]
pub fn create_quote(handler: State<DbConn>, channel: String, quote: Json<PostQuote>) -> JsonValue {
	let result = handler.lock().expect("db lock").create_quote(&channel, quote.into_inner());
	match result {
		Ok(id) => json!({
			"data": json!({
				"created": true,
				"id": id
			})
		}),
		Err(HandlerError::Error(e)) => generate_error(401, Some(e)),
		Err(e) => {
			println!("Internal error creating quote: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[get("/<channel>/random", rank = 1)]
pub fn get_random_quote(handler: State<DbConn>, channel: String) -> JsonValue {
	let quote = handler.lock().expect("db lock").get_random_quote(&channel);
	match quote {
		Ok(quote) => json!({ "data": quote }),
		Err(HandlerError::Error(_)) => json!({ "data": json!({}) }),
		Err(e) => {
			println!("Internal error getting quote: {:?}", e);
			generate_error(500, None)
		}
	}
}

#[delete("/<channel>/<id>", rank = 4)]
pub fn delete_quote(handler: State<DbConn>, channel: String, id: u32) -> JsonValue {
	let result = handler.lock().expect("db lock").delete_quote(&channel, id);
	match result {
		Ok(()) => json!({ "data": json!({ "deleted": true }) }),
		Err(HandlerError::Error(_)) => json!({ "data": json!({}) }),
		Err(e) => {
			println!("Internal error getting quote: {:?}", e);
			generate_error(500, None)
		}
	}
}
