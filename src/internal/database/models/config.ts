
import { Model, Column } from "pims";

@Model({
	table: "configs"
})
export class Configs {

	@Column({ primary: true })
	public id: string;

	@Column() public channel: string;
	@Column() public repeat: RepeatConfig;
	@Column() public events: EventsConfig;
	@Column() public spam: SpamConfig;
}
