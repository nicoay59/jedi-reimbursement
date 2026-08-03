export function formatRupiah(value) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0
  }).format(Number(value || 0));
}

export function formatDate(value) {
  if (!value) return "-";

  const date = new Date(`${value}T00:00:00`);
  if (Number.isNaN(date.getTime())) return value;

  return new Intl.DateTimeFormat("id-ID", {
    day: "2-digit",
    month: "long",
    year: "numeric"
  }).format(date);
}


export function formatDateRange(startDate, endDate) {
  if (!startDate) return "-";
  if (!endDate || startDate === endDate) {
    return formatDate(startDate);
  }

  const start = new Date(`${startDate}T00:00:00`);
  const end = new Date(`${endDate}T00:00:00`);
  if (Number.isNaN(start.getTime()) || Number.isNaN(end.getTime())) {
    return `${startDate} – ${endDate}`;
  }

  if (
    start.getFullYear() === end.getFullYear() &&
    start.getMonth() === end.getMonth()
  ) {
    const startDay = new Intl.DateTimeFormat("id-ID", {
      day: "2-digit"
    }).format(start);
    const endText = new Intl.DateTimeFormat("id-ID", {
      day: "2-digit",
      month: "long",
      year: "numeric"
    }).format(end);
    return `${startDay}–${endText}`;
  }

  return `${formatDate(startDate)} – ${formatDate(endDate)}`;
}

export function formatDateTime(value) {
  if (!value) return "-";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;

  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
    timeStyle: "short"
  }).format(date);
}

export function formatFileSize(bytes) {
  const value = Number(bytes || 0);

  if (value < 1024) return `${value} B`;
  if (value < 1024 * 1024) {
    return `${(value / 1024).toFixed(1)} KB`;
  }

  return `${(value / (1024 * 1024)).toFixed(1)} MB`;
}

export function localISODate() {
  const now = new Date();
  const offset = now.getTimezoneOffset();
  return new Date(now.getTime() - offset * 60_000)
    .toISOString()
    .slice(0, 10);
}


export function formatDuration(hours) {
  const totalMinutes = Math.round(Number(hours || 0) * 60);
  const hourPart = Math.floor(totalMinutes / 60);
  const minutePart = totalMinutes % 60;

  if (hourPart === 0) {
    return `${minutePart} menit`;
  }

  if (minutePart === 0) {
    return `${hourPart} jam`;
  }

  return `${hourPart} jam ${minutePart} menit`;
}

export function calculateOvertimeDuration(startTime, endTime) {
  if (!startTime || !endTime) {
    return null;
  }

  const startMinutes = clockToMinutes(startTime);
  const endMinutes = clockToMinutes(endTime);

  if (startMinutes === null || endMinutes === null) {
    return null;
  }

  let durationMinutes = endMinutes - startMinutes;
  let crossesMidnight = false;

  if (durationMinutes <= 0) {
    durationMinutes += 24 * 60;
    crossesMidnight = true;
  }

  return {
    minutes: durationMinutes,
    hours: Math.round((durationMinutes / 60) * 100) / 100,
    label: formatDuration(durationMinutes / 60),
    crossesMidnight
  };
}

function clockToMinutes(value) {
  const match = /^([01]\d|2[0-3]):([0-5]\d)$/.exec(value);
  if (!match) {
    return null;
  }

  return Number(match[1]) * 60 + Number(match[2]);
}


export function currentMonthRange() {
  const now = new Date();
  const firstDay = new Date(
    now.getFullYear(),
    now.getMonth(),
    1
  );

  return {
    startDate: toLocalISODate(firstDay),
    endDate: toLocalISODate(now)
  };
}

export function toLocalISODate(value) {
  const date = new Date(value);
  const offset = date.getTimezoneOffset();

  return new Date(date.getTime() - offset * 60_000)
    .toISOString()
    .slice(0, 10);
}

export function shortDate(value) {
  if (!value) return "-";

  const date = new Date(`${value}T00:00:00`);
  if (Number.isNaN(date.getTime())) return value;

  return new Intl.DateTimeFormat("id-ID", {
    day: "2-digit",
    month: "short"
  }).format(date);
}
