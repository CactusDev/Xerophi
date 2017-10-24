
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { ConfigController } from "./controller";

export class ConfigRoute extends AbstractEndpoint {

	private controller = new ConfigController();

	public async init() {
		this.web.instance.route({
			method: "GET",
			path: "/channel/{channel}/config",
			config: {
				handler: (request, reply) => this.controller.getConfig(request, reply),
				auth: false
			}
		});
	}
}