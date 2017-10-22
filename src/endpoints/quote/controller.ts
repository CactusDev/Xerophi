
import * as Hapi from "hapi";
import * as Boom from "boom";

import { MongoHandler } from "../../mongo";

const moment = require("moment-strftime");

export class QuoteController {

	constructor(private mongo: MongoHandler) {

	}

	private async verifyData(data: string): Promise<boolean> {
		return true; // TODO: How does this even get verified?
	}

	public async getQuote(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const id = +request.params["id"];
		const channel = request.params["channel"];

		if (id) {
			const random = !!request.payload && !!request.payload.random || false;
			const response = await this.mongo.getQuote(channel, random, id);
			if (!response || response.deletedAt) {
				return reply(Boom.notFound());
			}
			delete response.deletedAt;
			return reply(response);
		} else {
			const response = await this.mongo.getAllQuotes(channel);
			for (let quote of response) {
				delete quote.deletedAt;
			}
			return reply(response);
		}
	}

	public async createQuote(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		if (!await this.verifyData(request.payload.quote)) {
			// Invalid data, tell the user.
			return reply(Boom.badData("Invalid quote data"));
		}

		const channel = request.params["channel"];
		const quoted = request.payload.quoted;

		const quote: Quote = {
			quoteId: -1,
			channel: channel,
			quoted: quoted,
			createdAt: moment().strftime("%a %b %d %H:%M:%S %Y"),
			deletedAt: null,
			enabled: true,
			count: 0,
			quote: request.payload.quote
		};
		await this.mongo.createQuote(quote)

		delete quote.deletedAt;
		reply(quote);
	}
}