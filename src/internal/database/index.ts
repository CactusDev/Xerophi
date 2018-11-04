
import { RethinkAdapter, RethinkAdapterOptions } from "pims-rethinkdb";

import { Aliases, Channels, Commands, Configs,
	     Permits, Points, Quotes, Repeats, Users } from "./models";

export class DatabaseHandler {

	private instance: RethinkAdapter;

	constructor(options: RethinkAdapterOptions) {
		options.models = [
			Aliases, Channels, Commands, Configs, Permits, Points,
			Quotes, Repeats, Users
		];

		this.instance = new RethinkAdapter(options);
	}

	public async setup() {
		await this.instance.ensure();
	}

	public async findCommandByName(channel: string, name: string): Promise<Commands> {
		return await this.instance.findOne(Commands, { channel, name });
	}

	public async findQuoteById(channel: string, quoteId: number): Promise<Quotes> {
		return await this.instance.findOne(Quotes, { channel, quoteId });
	}

	public async getChannel(token: string): Promise<Channels> {
		return await this.instance.findOne(Channels, { token });
	}

	public async getRepeatsForChannel(channel: string): Promise<Repeats[]> {
		return await this.instance.find(Repeats, { channel });
	}

	public async getConfigForChannel(channel: string): Promise<Configs> {
		return await this.instance.findOne(Configs, { channel });
	}
}
