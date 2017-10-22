
const jwt = require("jsonwebtoken");

interface Authorized {
	[key: string]: {
		exp: number;
		id: string;
		valid: boolean;
	}
}

const data: Authorized = {};

export class Authorization {
	
	public static async give(session: any, secret: string): Promise<string> {
		// TODO: Put the generated key into redis instead
		const key = jwt.sign(session, secret);
		data[key] = session;
		return key;
	}

	public static async isValid(key: string): Promise<boolean> {
		return !!data[key] && data[key].exp > new Date().getDate();
	}

	public static async remove(key: string) {
		delete data[key];
	}
}