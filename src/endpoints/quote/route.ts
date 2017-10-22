
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { QuoteController } from "./controller";

export class QuoteRoute extends AbstractEndpoint {

	private controller = new QuoteController();

	public async init() {
		this.web.instance.route({
			method: "GET",
			path: "/{channel}/quote/{id}",
			config: {
				handler: (request, reply) => this.controller.getQuote(request, reply),
				auth: false
			}
		});
	}
}