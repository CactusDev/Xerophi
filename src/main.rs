#![feature(proc_macro_hygiene, decl_macro, result_map_or_else)]

mod endpoints;
mod database;

#[macro_use]
extern crate rocket;
#[macro_use]
extern crate rocket_contrib;

extern crate mongodb;
extern crate bson;
extern crate chrono;

extern crate serde;
extern crate serde_json;
#[macro_use]
extern crate serde_derive;
extern crate argon2;

use std::sync::Mutex;

pub type DbConn<'cfg> = Mutex<crate::database::handler::DatabaseHandler<'cfg>>;

fn main() {
	let mut connection = database::handler::DatabaseHandler::new("192.168.99.100", 32769, "123123123123", "123123123123");
	match connection.connect("cactus", "c", "c") {
		Ok(()) => println!("Connected!"),
		Err(e) => println!("Error: {:?}", e)
	};

    rocket::ignite()
    	.manage(Mutex::new(connection))
	    .mount("/channel", routes! [
	    	endpoints::channel::get_channel, endpoints::channel::create_channel
	    ])
	    .mount("/state", routes! [
	    	endpoints::state::get_channel_state, endpoints::state::get_channel_service_state
	    ])
	    .mount("/auth", routes! [
	    	endpoints::authorization::get_service_auth, endpoints::authorization::update_service_auth
	    ])
	    .mount("/quote", routes! [
	    	endpoints::quote::get_quote, endpoints::quote::get_random_quote,
	    	endpoints::quote::get_quote_by_id, endpoints::quote::create_quote,
	    	endpoints::quote::delete_quote
	    ])
	    .mount("/command", routes! [
	    	endpoints::command::get_commands, endpoints::command::get_command,
	    	endpoints::command::create_command, endpoints::command::delete_command,
	    	endpoints::command::edit_command, endpoints::command::update_count
	    ])
	    .mount("/trust", routes! [
	    	endpoints::trusts::get_trust, endpoints::trusts::delete_trust,
	    	endpoints::trusts::create_trust
	    ])
	    .mount("/alias", routes! [
	    	endpoints::alias::get_aliases, endpoints::alias::get_alias,
	    	endpoints::alias::create_alias, endpoints::alias::delete_alias
	    ])
	    .launch();
}
