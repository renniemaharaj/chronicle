export function notify(target: { toast: string }, message: string) {
  target.toast = message;
  setTimeout(() => {
    if (target.toast === message) {
      target.toast = '';
    }
  }, 3500);
}
