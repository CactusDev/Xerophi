
import { Config } from "./config";
import * as Hapi from "hapi";
import { Injectable } from "@angular/core";
import { AbstractEndpoint, CommandRoute } from "./endpoints";

@Injectable()
export class Web {
	private _instance: Hapi.Server;

	private endpoints: AbstractEndpoint[] = [];

	constructor(protected config: Config) {

	}

	public async start() {
		console.log("Starting...");

		this._instance = new Hapi.Server();

		this._instance.connection({
			port: this.config.web.port,
			routes: {
				cors: true
			}
		});
		this._instance.on("request-error", (req: any, error: any) => console.error(error));
		console.log("Creating endpoints...");
		this.endpoints.push(new CommandRoute(this, this.config));
		console.log("Done!");

		console.log("Initializing endpoints...");
		this.endpoints.forEach(async router => await router.init());
		console.log(`Done! Initialized ${this.endpoints.length} endpoints!`);
		this._instance.start();
		console.log(`Ready! :${this.config.web.port}`);
	}

	public get instance() {
		return this._instance;
	}
}