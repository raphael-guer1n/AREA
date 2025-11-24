import { Injectable, BadRequestException } from '@nestjs/common';
import { User } from './user.entity';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';

@Injectable()
export class UsersService {
  private users: User[] = [];
  private counter = 1;

  constructor() {
    const defaultUser: User = {
      id: this.counter,
      name: 'Name',
      email: 'Email',
      password: 'Password',
    };
    this.users.push(defaultUser);
  }

  findAll(): User[] {
    return this.users;
  }

  create(dto: CreateUserDto): string {
    this.counter += 1;
    const newUser: User = {
      id: this.counter,
      name: dto.name,
      email: dto.email,
      password: dto.password,
    };
    this.users.push(newUser);
    return 'User added successfully';
  }

  delete(id: number): string {
    const index = this.users.findIndex((u) => u.id === id);
    if (index === -1) {
      throw new BadRequestException('User does not exist');
    }
    this.users.splice(index, 1);
    return 'User deleted successfully';
  }

  update(id: number, dto: UpdateUserDto): string {
    if (dto.id !== id) {
      throw new BadRequestException('Id does not match');
    }
    const index = this.users.findIndex((u) => u.id === id);
    if (index === -1) {
      throw new BadRequestException('User does not exist');
    }
    this.users[index] = { ...dto };
    return 'User updated successfully';
  }
}
