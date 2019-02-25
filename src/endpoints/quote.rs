
use rocket_contrib::json::Json;

use rocket::{
	State, Response,
	http::Status
};

use crate::{
	DbConn, endpoints::{generate_error, generate_response},
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
pub fn get_quote<'r>(handler: State<DbConn>, channel: String) -> Response<'r> {
	let quotes = handler.lock().expect("db lock").get_quote(&channel, None);
	match quotes {
		Ok(quotes) => generate_response(Status::Ok, json!({ "data": quotes })),
		Err(HandlerError::Error(_)) => generate_response(Status::NotFound, json!({ "data": json!([]) })),
		Err(e) => {
			println!("Internal error getting quote: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[get("/<channel>/<id>", rank = 3)]
pub fn get_quote_by_id<'r>(handler: State<DbConn>, channel: String, id: u32) -> Response<'r> {
	let quote = handler.lock().expect("db lock").get_quote(&channel, Some(id));
	match quote {
		Ok(quote) => generate_response(Status::Ok, json!({ "data": quote[0] })),
		Err(HandlerError::Error(_)) => generate_response(Status::NotFound, json!({ "data": json!([]) })),
		Err(e) => {
			println!("Internal error getting quote: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[post("/<channel>/create", rank = 2, format = "json", data = "<quote>")]
pub fn create_quote<'r>(handler: State<DbConn>, channel: String, quote: Json<PostQuote>) -> Response<'r> {
	let result = handler.lock().expect("db lock").create_quote(&channel, quote.into_inner());
	match result {
		Ok(id) => generate_response(Status::Created, json!({
			"data": json!({
				"created": true,
				"id": id
			})
		})),
		Err(HandlerError::Error(_)) => generate_response(Status::BadRequest, json!({ "data": json!([]) })),
		Err(e) => {
			println!("Internal error creating quote: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[get("/<channel>/random", rank = 1)]
pub fn get_random_quote<'r>(handler: State<DbConn>, channel: String) -> Response<'r> {
	let quote = handler.lock().expect("db lock").get_random_quote(&channel);
	match quote {
		Ok(quote) => generate_response(Status::Ok, json!({ "data": quote })),
		Err(HandlerError::Error(_)) => generate_response(Status::NotFound, json!({ "data": json!([]) })),
		Err(e) => {
			println!("Internal error getting quote: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[delete("/<channel>/<id>", rank = 4)]
pub fn delete_quote<'r>(handler: State<DbConn>, channel: String, id: u32) -> Response<'r> {
	let result = handler.lock().expect("db lock").delete_quote(&channel, id);
	match result {
		Ok(()) => generate_response(Status::NoContent, json!({})),
		Err(HandlerError::Error(_)) => generate_response(Status::NotFound, json!({"data": {}})),
		Err(e) => {
			println!("Internal error getting quote: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}
