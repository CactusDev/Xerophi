
type Role = "banned" | "user" | "subscriber" | "moderator" | "owner";

interface Component {
    type: "text" | "emoji" | "tag" | "url" | "variable";
    data: string;
}

/**
 * Text {@see Component}, containing raw text.
 *
 * @interface Text
 */
interface TextComponent extends Component {
    type: "text";
}

/**
 * Emoji {@see Component}, containing an emoji.
 * If the emoji is a standard Unicode emoji, its alpha code should be used.
 * Otherwise, a consistent name should be chosen, prefixed with a period.
 *
 * @interface Emoji
 */
interface EmojiComponent extends Component {
    type: "emoji";
}

/**
 * Tag {@see Component}, containing a tag. No prefix, such as an @ symbol,
 * should be stored.
 *
 * @interface Tag
 */
interface TagComponent extends Component {
    type: "tag";
}

/**
 * URL {@see Component}, containing a URL.
 *
 * @interface URL
 */
interface URLComponent extends Component {
    type: "url";
}

/**
 * Variable {@see Component}, containing isolated variable data.
 *
 * @example {
 *     type: "variable",
 *     data: "ARG3|reverse|title"
 * }
 *
 * @example {
 *     type: "variable",
 *     data: "USER"
 *  }
 *
 * @interface Variable
 */
interface VariableComponent extends Component {
    type: "variable";
}

/**
 * Message packet
 *
 * @interface CactusMessagePacket
 */
interface CactusMessagePacket {
    type: "message";
    text: Component[];
    action: boolean;
}

interface Command {
	name: string;
	channel: string;
	response: CactusMessagePacket[];
	count: number;
	enabled: boolean;
	restrictions: {
		service: string[];
		role: Role;
	}
}

interface Quote {
	quoteId: number;
	channel: string;
	quoted: string;
	when: string;
	enabled: boolean;
	count: number;
	quote: Component[];
}

interface SpamConfig<T> {
    action: 'ignore' | 'purge' | 'timeout' | 'ban';
    value: T;
    warnings: number;
}

interface EventConfig {
    message: string;
    enabled: boolean;
}

interface Config {
    repeat: {
  	  disabled: boolean;
      onlyLive: boolean;
      defaultMinimum: number;
    };
  
    events: {
        follow: EventConfig;
        subscribe: EventConfig;
        host: EventConfig;
        join: EventConfig;
        leave: EventConfig;
    };
  
    whitelistedURLs: string[];
  
    spam: {
        allowUrls: SpamConfig<boolean>;
        maxCaps: SpamConfig<number>;
        maxEmoji: SpamConfig<number>;

	    keywords: {
		    blacklist: string[];
	        whitelist: string[];
	    };
    }
}

interface Repeat {
	text: string;
	isCommand: boolean;
	delay: number;
	disabled: boolean
}

interface ServiceAuth {
	accountName: string;
	access: string;
	refresh?: string;
	expires?: string;
}

interface Chatters {
	[name: string]: {
		points: number;
	}
}

interface Channel {
	repeats: Repeat[];
	username: string;
	service: string;
	uuid: string;
	trusts: string[];
	permits: string[];
	chatters: Chatters;
	config: Config;
}

interface User {
	username: string;
	deletedAt: string;
	uuid: string;
	passwordHash: string;
	channels: Channel[];
	scopes: string[];
	commands: Command[];
}