
interface ServiceAuthorization {
	service: string,
	token: string,
	authKey: string,
	refresh?: string,
	expires?: number
}

interface BaseSpamConfig<T> {
  action: 'ignore' | 'purge' | 'timeout' | 'ban';
  value: T;
  warnings: number;
}

interface EventConfig {
  message: string;
  enabled: boolean;
}

interface RepeatConfig {
	disabled: boolean;
	onlyLive: boolean;
	defaultMinimum: number;
}

interface EventsConfig {
	follow: EventConfig,
	subscribe: EventConfig,
	host: EventConfig,
	join: EventConfig,
	leave: EventConfig
}

interface KeywordConfig {
	blacklist: string[];
	whitelist: string[];
}

interface SpamConfig {
	allowUrls: BaseSpamConfig<boolean>,
	maxCapsScore: BaseSpamConfig<number>,
	maxEmoji: BaseSpamConfig<number>,

	keywords: KeywordConfig,
	whitelistedUrls: string[]
}

interface Component {
	type: "text" | "emoji" | "tag" | "url",
	data: string
}

interface CommandResponse {
	text: Component[],
	action: boolean,
	target?: string,
	role: number,
	user: string
}
