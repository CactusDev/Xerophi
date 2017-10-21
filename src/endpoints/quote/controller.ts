
import * as Hapi from "hapi";
import * as Boom from "boom";

export class QuoteController {

	public async getQuote(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const id = +request.params["id"];
		const channel = request.params["channel"];

		// TODO: Make this actually pull from a database and display information
		const response: Quote = {
			quoteId: id,
			channel: channel,
			quoted: "2Cubed",
			when: "2017 10-21 2:14",
			enabled: true,
			count: 0,
			quote: [
				{
					type: "text",
					data: "I will hit you with a potato"
				},
				{
					type: "emoji",
					data: "green_heart"
				}
			],
		};
		reply(response);
	}
}