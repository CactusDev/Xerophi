
import { Config } from "../config";
import { theTime } from "../util";

import * as Mongo from "mongodb";

export class MongoHandler {

	private connection: Mongo.Db;
	private commands: Mongo.Collection;
	private quotes: Mongo.Collection;

	constructor(private config: Config) {
	}

	public async connect() {
		const connectionURL = `mongodb://${this.config.mongo.username}:${this.config.mongo.password}@${this.config.mongo.host}:${this.config.mongo.port}/${this.config.mongo.database}`;
		Mongo.MongoClient.connect(connectionURL, {
				authSource: this.config.mongo.authdb
			},
			(err, db) => {
			if (err) {
				console.error(err);
				return;
			}
			this.connection = db;

			this.quotes = this.connection.collection("quotes");
			this.commands = this.connection.collection("commands");
		});
	}

	public async createQuote(quote: Quote) {
		const recent = await this.quotes.find({ channel: quote.channel }).sort({ quoteId: -1 }).limit(1).toArray();
		const quoteId = recent.length == 0 ? 1 : recent[0].quoteId + 1
		quote.quoteId = quoteId;

		this.quotes.insertOne(quote);
	}

	public async getAllQuotes(channel: string): Promise<Quote[]> {
		return await this.quotes.find({ channel }).toArray();
	}

	public async getQuote(channel: string, random: boolean, quoteId: number): Promise<Quote> {
		if (!quoteId && random) {
			const quotes = await this.quotes.aggregate([
				{
					"$sample": { "size": 1 }
				},
				{
					"$match": { channel }
				}
			]).toArray();
			if (quotes.length === 0) {
				return null;
			}
			return quotes[0];
		}
		if (quoteId === -1) {
			return null;
		}
		return await this.quotes.findOne({ channel, quoteId });
	}

	public async getCommand(channel: string, name: string): Promise<Command> {
		const commands = await this.commands.find({ channel, name }).toArray();
		return commands.length == 0 ? null : commands[0];
	}

	public async createCommand(channel: string, name: string, response: CactusMessagePacket[], role: string): Promise<boolean> {
		if (await this.getCommand(channel, name)) {
			return false;
		}
		const command: Command = {
			name: name,
			channel,
			response,
			deletedAt: null,
			count: 0,
			enabled: true,
			restrictions: {
				service: [],
				role
			}
		};
		this.commands.insertOne(command);
		return true;
	}

	public async editCommandAttribute(attribute: string, value: any, command: string, channel: string): Promise<boolean> {
		// Make sure the command exists
		const dbCommand: any = await this.getCommand(channel, command);
		if (!dbCommand) {
			return false;
		}
		// Check if it's the special types
		if (attribute === "role" || attribute === "service") {
			dbCommand.restrictions[attribute] = value;
			const result = await this.commands.updateOne({ channel, name: command }, dbCommand);
			return result.matchedCount === 1;
		} else if (attribute === "count") {
			// Count is *super* special. We need to figure out what we're doing
			const prefix = (<string>value).substring(0, 1);
			const valueIsNum = !!+value;
			const isNum = !!+prefix && valueIsNum;
			const remaining = (<string>value).substring(1);
			const remainingIsNum = !!+remaining;

			if (isNum) {
				// Setting
				if (!valueIsNum) {
					return true;
				}
				dbCommand.count = +value
				const result = await this.commands.updateOne({ channel, name: command }, dbCommand);
				return result.matchedCount === 1;
			} else if (remainingIsNum) {
				if (prefix === "+") {
					// Adding to the count
					if (!valueIsNum) {
						return true;
					}
					dbCommand.count += +remaining
					const result = await this.commands.updateOne({ channel, name: command }, dbCommand);
					return result.matchedCount === 1;
				} else if (prefix === "-") {
					// Subtracting
					if (!valueIsNum) {
						return true;
					}
					dbCommand.count -= +remaining
					const result = await this.commands.updateOne({ channel, name: command }, dbCommand);
					return result.matchedCount === 1;
				}
			}
		}
		// Update the attribute
		dbCommand[attribute] = value;
		const result = await this.commands.updateOne({ channel, name: command }, dbCommand);
		return result.matchedCount === 1;
	}

	public async softDeleteCommand(name: string, channel: string): Promise<boolean> {
		const command = await this.getCommand(channel, name);
		if (!command) {
			return false;
		}
		command.enabled = false;
		command.deletedAt = theTime();
		const result = await this.commands.updateOne({ channel, name }, command);
		return result.matchedCount === 1;
	}

	public async softDeleteQuote(id: number, channel: string): Promise<boolean> {
		const quote = await this.getQuote(channel, false, id);
		if (!quote || quote.deletedAt) {
			return false;
		}
		quote.deletedAt = theTime();
		quote.enabled = false;
		const result = await this.quotes.updateOne({ channel, quoteId: id }, quote);
		return result.matchedCount === 1;
	}
}