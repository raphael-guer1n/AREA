import { Controller, Get, Post, Delete, Put, Param, Body } from '@nestjs/common';
import { UsersService } from './users.service';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';

@Controller('users')
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @Get()
  getUsers() {
    const data = this.usersService.findAll();
    return { data };
  }

  @Post()
  addUser(@Body() body: CreateUserDto) {
    const message = this.usersService.create(body);
    return { data: message };
  }

  @Delete(':id')
  deleteUser(@Param('id') id: string) {
    const message = this.usersService.delete(Number(id));
    return { data: message };
  }

  @Put(':id')
  updateUser(@Param('id') id: string, @Body() body: UpdateUserDto) {
    const message = this.usersService.update(Number(id), body);
    return { data: message };
  }
}
