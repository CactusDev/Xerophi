
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { LoginController } from "./controller";

export class LoginRoute extends AbstractEndpoint {

	private controller = new LoginController();

	public async init() {
		this.web.instance.route({
			method: "GET",
			path: "/login",
			config: {
				handler: (request, reply) => this.controller.attemptLogin(request, reply, this.config.authentication.secret),
				auth: false
			}
		});
	}
}