
import * as Hapi from "hapi";
import * as Boom from "boom";

const moment = require("moment-strftime");

export class QuoteController {

	private async verifyData(data: string): Promise<boolean> {
		return true; // TODO: How does this even get verified?
	}

	public async getQuote(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const id = +request.params["id"];
		const channel = request.params["channel"];

		// TODO: Make this actually pull from a database and display information
		const response: Quote = {
			quoteId: id,
			channel: channel,
			quoted: "2Cubed",
			createdAt: "2017 10-21 2:14",
			deletedAt: null,
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
			]
		};
		delete response.deletedAt;
		reply(response);
	}

	public async createQuote(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const id = 10; // TODO: This should be the proper id after inserting
		if (!await this.verifyData(request.payload.quote)) {
			// Invalid data, tell the user.
			return reply(Boom.badData("Invalid quote data"));
		}

		const channel = request.params["channel"];
		const quoted = request.payload.quoted;

		const quote: Quote = {
			quoteId: id,
			channel: channel,
			quoted: quoted,
			createdAt: moment().strftime("%a %b %d %H:%M:%S %Y"),
			deletedAt: null,
			enabled: true,
			count: 0,
			quote: request.payload.quote
		};
		// TODO: Insert
		delete quote.deletedAt;
		reply(quote);
	}
}