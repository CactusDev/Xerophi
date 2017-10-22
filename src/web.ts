
import { Config } from "./config";
import * as Hapi from "hapi";
import { Injectable } from "@angular/core";
import { AbstractEndpoint, CommandRoute, QuoteRoute, ChannelRoute, ConfigRoute, LoginRoute } from "./endpoints";

import { Authorization } from "./authorization";
import { MongoHandler } from "./mongo";

@Injectable()
export class Web {
	private _instance: Hapi.Server;

	private endpoints: AbstractEndpoint[] = [];
	public mongo: MongoHandler;

	constructor(protected config: Config) {

	}

	public async start(mongo: MongoHandler) {
		this.mongo = mongo;

		console.log("Starting...");

		const validate = (decoded: any, request: Hapi.Request, callback: any) => {
			// TODO: Create a better check here to make sure the user exists
			Authorization.isValid(request.headers.authorization).then(valid => callback(null, valid, { scope: decoded.scopes }));
		};

		this._instance = new Hapi.Server();

		this._instance.on("response", (request) =>
			console.log(`${request.info.remoteAddress}: ${request.method.toUpperCase()} <${request.response.statusCode}> ${request.url.path}`));

		this._instance.connection({
			port: this.config.web.port,
			routes: {
				cors: true
			}
		});

		this._instance.on("request-error", (req: any, error: any) => console.error(error));
	
		this._instance.register(require("hapi-auth-jwt2"), (err) => {
			if (err) {
				throw err;
			}

			this._instance.auth.strategy("jwt", "jwt", {
				key: this.config.authentication.secret,
				validateFunc: validate,
				verifyOptions: {
					algorithms: ["HS256"]
				}
			});

			this._instance.auth.default("jwt");

			this._instance.route([
				{
					method: "GET",
					path: "/",
					config: {
						auth: false
					},
					handler: (request, reply) => reply({a: "thing"})
				},
				{
					method: "GET",
					path: "/a",
					config: {
						auth: "jwt"
					},
					handler: (request, reply) => reply({b: "stuff"})
				}
			]);

			console.log("Creating endpoints...");
			this.endpoints.push(new CommandRoute(this, this.config), new QuoteRoute(this, this.config),
				new ChannelRoute(this, this.config), new ConfigRoute(this, this.config), new LoginRoute(this, this.config));
			console.log("Done!");

			console.log("Initializing endpoints...");
			this.endpoints.forEach(async router => await router.init());
			console.log(`Done! Initialized ${this.endpoints.length} endpoints!`);
			this._instance.start();
			console.log(`Ready! :${this.config.web.port}`);
		});
	}

	public get instance() {
		return this._instance;
	}
}