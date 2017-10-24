
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { UserController } from "./controller";

export class UserRoute extends AbstractEndpoint {

	private controller = new UserController();

	public async init() {
		this.web.instance.route({
			method: "GET",
			path: "/users/login",
			config: {
				handler: (request, reply) => this.controller.attemptLogin(request, reply, this.config.authentication.secret),
				auth: false
			}
		});
	}
}