# Migrations

This directory contains database migration scripts for the project. These scripts are used to manage changes to the database schema over time, ensuring that the database structure remains consistent and up-to-date.

## Usage

- To create a new migration, use the following command:

```sh

make migrate/create name=<migration_name>
```

- To apply the migrations, use the following command:

```sh
make migrate/up
```
- To revert the migrations, use the following command:

```sh
make migrate/down
```
