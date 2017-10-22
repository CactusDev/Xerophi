
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
		const random = !!request.query && !!request.query.random || false;

		if (id) {
			const response = await this.mongo.getQuote(channel, false, id);
			if (!response || response.deletedAt) {
				return reply(Boom.notFound("Invalid quote"));
			}
			delete response.deletedAt;
			if (!response.deletedAt && response.enabled) {
				return reply(response);
			}
			return reply(Boom.notFound("Invalid quote"));
		} else {
			if (random) {
				const response = await this.mongo.getQuote(channel, true);
				if (!response || response.deletedAt) {
					return reply(Boom.notFound("Invalid quote"));
				}
				delete response.deletedAt;
				if (!response.deletedAt && response.enabled) {
					return reply(response);
				}
				return reply(Boom.internal("Got invalid quote from filter?"));
			}
 			const response = await this.mongo.getAllQuotes(channel);
			for (let quote of response) {
				if (quote.deletedAt) {
					const index = response.indexOf(quote);
					if (index !== -1) {
						response.splice(index, 1);
					}
				} else {
					delete quote.deletedAt;					
				}
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
		return reply({
			created: true,
			name
		}).code(201);
	}

	public async deleteQuote(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params.channel;
		const quoteId = +request.params.id;
		if (!quoteId) {
			return reply(Boom.badData("Quote id must be a number."));
		}

		const deleted = await this.mongo.softDeleteQuote(quoteId, channel);
		if (!deleted) {
			return reply(Boom.notFound("Invalid quote"));
		}
		return reply({}).code(204);
	}
}