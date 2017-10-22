
const moment = require("moment-strftime");
const argon2 = require("argon2");

const validRoles = ["banned", "user", "subscriber", "moderator", "owner"];

export function isValidRole(role: string): boolean {
	return validRoles.indexOf(role) > -1;
}

export function theTime(): string {
	return moment().strftime("%a %b %d %H:%M:%S %Y");
}

export async function compare(password: string, hash: string): Promise<boolean> {
	try {
		return await argon2.verify(hash, password);
	} catch (e) {
		return false;
	}
}

export async function hash(password: string): Promise<string> {
	try {
		return await argon2.hash(password);
	} catch (e) {
		return null;
	}
}
