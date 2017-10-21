
export class Config {
    public mongo: {
        host: string;
        port: number;
        username: string;
        password: string;
        database: string;
    }

    public web: {
        port: number;
    }

    public authentication: {
        secret: string;
    }
}