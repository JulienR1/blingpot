export const timestampFormatter = new Intl.DateTimeFormat("fr-CA", {
  day: "numeric",
  month: "short",
});
export const moneyFormatter = new Intl.NumberFormat("fr-CA", {
  style: "currency",
  currency: "CAD",
});
