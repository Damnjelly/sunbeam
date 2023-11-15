import { Param } from "./action.js";
export type Manifest = {
    title: string;
    platforms?: Platform[];
    description?: string;
    requirements?: Requirement[];
    items?: RootItem[];
    commands: CommandSpec[];
};
type RootItem = {
    command: string;
    title?: string;
    params?: Record<string, Param>;
};
type Platform = "linux" | "macos";
type Requirement = {
    name: string;
    link?: string;
};
export type CommandSpec = {
    name: string;
    title: string;
    mode: "list" | "detail" | "tty" | "silent";
    hidden?: boolean;
    description?: string;
    params?: Input[];
};
type PayloadParams = Record<string, string | boolean | number>;
export type Payload<T extends PayloadParams = PayloadParams, V extends Record<string, any> = Record<string, any>> = {
    command: string;
    params: T;
    preferences: V;
    query?: string;
    cwd: string;
};
type InputProps = {
    name: string;
    required: boolean;
};
type TextField = InputProps & {
    type: "text";
    title: string;
    defaut?: string;
    placeholder?: string;
};
type NumberField = InputProps & {
    type: "number";
    title: string;
    default?: number;
    placeholder?: string;
};
type TextArea = InputProps & {
    type: "textarea";
    title: string;
    defaut?: string;
    placeholder?: string;
};
type Password = InputProps & {
    type: "password";
    title: string;
    defaut?: string;
    placeholder?: string;
};
type Checkbox = InputProps & {
    type: "checkbox";
    label: string;
    title?: string;
    defaut?: boolean;
};
export type Input = TextField | TextArea | Password | Checkbox | NumberField;
export {};
