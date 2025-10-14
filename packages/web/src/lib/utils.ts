import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

type Indexable<T> = {
  [K in keyof T]: T[K] extends string | number | symbol ? T[K] : never;
};

type Always<T> = { [K in keyof T as T[K] extends never ? never : K]: T[K] };

export function dict<
  T extends Record<string, unknown>,
  K extends keyof Always<Indexable<T>>
>(elements: Array<T>, key: K): Record<Indexable<T>[typeof key], T> {
  const out = {} as Record<Indexable<T>[typeof key], T>;
  for (const element of elements) {
    const identifier = element[key as keyof T];
    out[identifier as Indexable<T>[typeof key]] = element;
  }
  return out;
}
