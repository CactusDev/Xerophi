
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
pub struct PostSocial {
	pub url: String
}

#[get("/<channel>")]
pub fn get_socials<'r>(handler: State<DbConn>, channel: String) -> Response<'r> {
	let socials = handler.lock().expect("db lock").get_socials(&channel);
	match socials {
		Ok(socials) => generate_response(Status::Ok, json!({ "data": socials })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error getting socials: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[get("/<channel>/<service>")]
pub fn get_social<'r>(handler: State<DbConn>, channel: String, service: String) -> Response<'r> {
	let social = handler.lock().expect("db lock").get_social_service(&channel, &service);
	match social {
		Ok(social) => generate_response(Status::Ok, json!({ "data": social })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error getting social: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}

#[patch("/<channel>/<service>", format = "json", data = "<social>")]
pub fn create_social<'r>(handler: State<DbConn>, channel: String, service: String, social: Json<PostSocial>) -> Response<'r> {
	let social = handler.lock().expect("db lock").create_social_service(&channel, &service, &social.into_inner().url);
	match social {
		Ok(social) => generate_response(Status::Ok, json!({ "data": social })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error creating social: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}


#[delete("/<channel>/<service>")]
pub fn delete_social<'r>(handler: State<DbConn>, channel: String, service: String) -> Response<'r> {
	let result = handler.lock().expect("db lock").remove_social(&channel, &service);
	match result {
		Ok(()) => generate_response(Status::Ok, json!({ "data": json!({ "deleted": true }) })),
		Err(HandlerError::Error(e)) => generate_response(Status::NotFound, generate_error(404, Some(e))),
		Err(e) => {
			println!("Internal error creating social: {:?}", e);
			generate_response(Status::InternalServerError, generate_error(500, None))
		}
	}
}
