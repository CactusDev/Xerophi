
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { QuoteController } from "./controller";

export class QuoteRoute extends AbstractEndpoint {

	private controller: QuoteController;

	public async init() {
		this.controller = new QuoteController(this.web.mongo);

		this.web.instance.route([
			{
				method: "GET",
				path: "/{channel}/quote/{id}",
				config: {
					handler: (request, reply) => this.controller.getQuote(request, reply),
					auth: false
				}
			},
			{
				method: "GET",
				path: "/{channel}/quote",
				config: {
					handler: (request, reply) => this.controller.getQuote(request. reply),
					auth: false
				}
			},
			{
				method: "POST",
				path: "/{channel}/quote",
				config: {
					handler: (request, reply) => this.controller.createQuote(request, reply),
					auth: {
						scope: ["user:quote:create"],
					}
				}
			},
			{
				method: "DELETE",
				path: "/{channel}/quote/{id}",
				config: {
					handler: (request, reply) => this.controller.deleteQuote(request, reply)
				}
			}
		]);
	}
}