
import * as Hapi from "hapi";
import * as Boom from "boom";
import { isValidRole } from "../../util";

import { MongoHandler } from "../../mongo";

const validEditable = ["enabled", "name", "response", "role", "service", "count"];

export class CommandController {

	constructor(private mongo: MongoHandler) {

	}

	public async validResponse(response: CactusMessagePacket[]): Promise<boolean> {
		return true; // TODO
	}

	public async getCommand(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const name = request.params["command"];
		const channel = request.params["channel"];

		const command = await this.mongo.getCommand(channel, name);
		if (!command) {
			return reply(Boom.notFound("Invalid command"));
		}
		return reply(command);
	}

	public async createCommand(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params["channel"];
		const name = request.params["command"];

		if (!request.payload || !request.payload.response || !request.payload.role) {
			return reply(Boom.badData("Must supply a payload with role & response"));
		}

		const response: CactusMessagePacket[] = !!request.payload && !!request.payload.response ? request.payload.response : null;
		if (!response || !await this.validResponse(response)) {
			return reply(Boom.badData("Invalid response"));
		}

		const role: string = request.payload.role;
		if (!isValidRole(role)) {
			return reply(Boom.badData("Invalid role"));
		}

		// Everything is valid, so lets make a command!
		const created = await this.mongo.createCommand(channel, name, response, role);
		if (!created) {
			return reply(Boom.conflict("Command already exists"));
		}
		// Created, display a nice message
		return reply({
			created: true,
			name
		}).code(201);
	}

	public async updateCommand(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params["channel"];
		const command = request.params["command"];
		
		for (let key of Object.keys(request.payload)) {
			const index = validEditable.indexOf(key);
			if (index == -1) {
				return reply(Boom.badData("Invalid attribute"));
			}

			const type = validEditable[index];
			const updated = await this.mongo.editCommandAttribute(key, request.payload[key], command, channel);

			if (!updated) {
				return reply(Boom.notFound("Invalid command."));
			}
		}
		return reply({}).code(204);
	}
}