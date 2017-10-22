
export const ADD_COMMAND = "user:command:create";
export const ADD_QUOTE = "user:command:create";
export const ADD_USER = "user:create";

export const REMOVE_COMMAND = "user:command:create";
export const REMOVE_QUOTE = "user:command:create";
export const REMOVE_USER = "user:command:create";

export const EDIT_COMMAND = "user:command:create";
export const EDIT_USER = "user:command:create";

export const fullRoles = [
	ADD_COMMAND,
	REMOVE_COMMAND,
	EDIT_COMMAND,
	ADD_QUOTE,
	REMOVE_QUOTE,
	ADD_USER,
	REMOVE_USER,
	EDIT_USER
];

// @TEMP: Make this the real list
export const defaultScopes = fullRoles;
