package backend

type Move struct{
	Id int64; //could be string if UUID
	Name string;
	Description string;
	Image string;
}