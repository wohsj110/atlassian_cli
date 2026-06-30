#!/usr/bin/env node

const fs = require("fs");
const os = require("os");
const path = require("path");

function usage() {
  console.log("Usage: atlassian-agent-skill install [--dest DIR]");
}

function copyDir(src, dest) {
  fs.mkdirSync(dest, { recursive: true });
  for (const entry of fs.readdirSync(src, { withFileTypes: true })) {
    const from = path.join(src, entry.name);
    const to = path.join(dest, entry.name);
    if (entry.isDirectory()) {
      copyDir(from, to);
    } else {
      fs.copyFileSync(from, to);
    }
  }
}

function defaultDest() {
  if (process.env.CODEX_HOME) {
    return path.join(process.env.CODEX_HOME, "skills");
  }
  return path.join(os.homedir(), ".codex", "skills");
}

function skillsRoot() {
  const candidates = [
    path.resolve(__dirname, "..", "skills"),
    path.resolve(__dirname, "..", "..", "..", "skills"),
  ];
  for (const candidate of candidates) {
    if (fs.existsSync(path.join(candidate, "Jira", "SKILL.md"))) {
      return candidate;
    }
  }
  throw new Error("Could not find bundled Jira and Confluence skills.");
}

function main(argv) {
  const command = argv[2];
  if (!command || command === "--help" || command === "-h") {
    usage();
    return 0;
  }
  if (command !== "install") {
    console.error(`Unknown command: ${command}`);
    usage();
    return 2;
  }

  let dest = defaultDest();
  for (let i = 3; i < argv.length; i += 1) {
    if (argv[i] === "--dest") {
      dest = argv[i + 1];
      i += 1;
    } else {
      console.error(`Unknown flag: ${argv[i]}`);
      return 2;
    }
  }

  const root = skillsRoot();
  for (const name of ["Jira", "Confluence"]) {
    copyDir(path.join(root, name), path.join(dest, name));
  }
  console.log(`Installed Jira and Confluence skills to ${dest}`);
  return 0;
}

process.exitCode = main(process.argv);
