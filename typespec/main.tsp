import "@typespec/http";
import "@typespec/openapi3";

using TypeSpec.Http;

@service(#{
  title: "User Management API",
})
@server("http://localhost:8080/api/v1", "Development server")
namespace UserManagementAPI;

/**
 * User model
 */
model User {
  /**
   * User ID (ULID format)
   */
  @pattern("^[0-9A-HJKMNP-TV-Z]{26}$")
  id: string;

  /**
   * User name
   */
  @minLength(1)
  @maxLength(100)
  name: string;

  /**
   * User email address
   */
  @format("email")
  email: string;

  /**
   * Creation timestamp
   */
  createdAt: utcDateTime;

  /**
   * Last update timestamp
   */
  updatedAt: utcDateTime;
}

/**
 * Create user request
 */
model CreateUserRequest {
  /**
   * User name
   */
  @minLength(1)
  @maxLength(100)
  name: string;

  /**
   * User email address
   */
  @format("email")
  email: string;
}

/**
 * Update user request
 */
model UpdateUserRequest {
  /**
   * User name
   */
  @minLength(1)
  @maxLength(100)
  name?: string;

  /**
   * User email address
   */
  @format("email")
  email?: string;
}

/**
 * User list response
 */
model UserList {
  /**
   * List of users
   */
  users: User[];

  /**
   * Total number of users
   */
  total: int32;
}

/**
 * Error response
 */
@error
model Error {
  /**
   * Error message
   */
  message: string;

  /**
   * Error code
   */
  code?: string;
}

@tag("users")
@route("/users")
interface Users {
  /**
   * Get all users
   */
  @get
  listUsers(
    /**
     * Maximum number of users to return
     */
    @query
    @minValue(1)
    @maxValue(100)
    limit?: int32 = 10,

    /**
     * Number of users to skip
     */
    @query
    @minValue(0)
    offset?: int32 = 0
  ): UserList | Error;

  /**
   * Create a new user
   */
  @post
  createUser(
    @body body: CreateUserRequest
  ): {
    @statusCode statusCode: 201;
    @body body: User;
  } | Error;

  /**
   * Get user by ID
   */
  @get
  @route("/{userId}")
  getUser(
    /**
     * User ID (ULID format)
     */
    @path
    @pattern("^[0-9A-HJKMNP-TV-Z]{26}$")
    userId: string
  ): User | Error;

  /**
   * Update user
   */
  @put
  @route("/{userId}")
  updateUser(
    /**
     * User ID (ULID format)
     */
    @path
    @pattern("^[0-9A-HJKMNP-TV-Z]{26}$")
    userId: string,

    @body body: UpdateUserRequest
  ): User | Error;

  /**
   * Delete user
   */
  @delete
  @route("/{userId}")
  deleteUser(
    /**
     * User ID (ULID format)
     */
    @path
    @pattern("^[0-9A-HJKMNP-TV-Z]{26}$")
    userId: string
  ): {
    @statusCode statusCode: 204;
  } | Error;
}
