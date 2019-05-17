/**
 * This file is part of the ego-cli distribution (https://github.com/egodigital/ego-cli).
 * Copyright (c) e.GO Digital GmbH, Aachen, Germany
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

/**
 * A command.
 */
export interface Command {
    /**
     * Executes the command.
     *
     * @param {CommandExecutionContext} context The context.
     *
     * @return {CommandExecuteResult|PromiseLike<CommandExecuteResult>} The result.
     */
    readonly execute: (context: CommandExecutionContext) => CommandExecuteResult | PromiseLike<CommandExecuteResult>;
}

/**
 * The result of an Command.#execute() invocation.
 */
export type CommandExecuteResult = void | null | undefined | number;

/**
 * Execution context of a command.
 */
export interface CommandExecutionContext {
    /**
     * The root directory of the command.
     */
    readonly root: string;
}


/**
 * A basic command.
 */
export abstract class CommandBase implements Command {
    /** @inheritdoc */
    public abstract execute(context: CommandExecutionContext): Promise<CommandExecuteResult>;
}
