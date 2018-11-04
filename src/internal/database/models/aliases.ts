
@Model({
	table: "aliases"
})
export class Aliases {

	@Column({ primary: true })
	public id: string;

	@Column() public channel: string;
	@Column() public alias: string;
	@Column() public to: string;
}
