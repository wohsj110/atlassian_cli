#!/usr/bin/env node

const fs = require("fs");
const os = require("os");
const path = require("path");
const { spawnSync } = require("child_process");

const SKILLS = ["atk-jira", "atk-cfl"];
const BINARIES = ["atk-jira", "atk-cfl"];
const REPO = "https://github.com/wohsj110/atlassian_cli";

function usage() {
  console.log(`Usage:
  skills add atlassian-agent [--target codex|claude|both] [--dest DIR] [--install-cli]
  atlassian-agent-skill install [--target codex|claude|both] [--dest DIR] [--install-cli]
  atlassian-agent-skill install-cli
  atlassian-agent-skill doctor

Examples:
  npx @wohsj110/skills add atlassian-agent
  npx atlassian-agent-skill install
  npx atlassian-agent-skill install --target codex --install-cli
  npx atlassian-agent-skill install --target claude
`);
}

function run(command, args, options = {}) {
  return spawnSync(command, args, {
    stdio: options.stdio || "pipe",
    encoding: "utf8",
  });
}

function hasCommand(command) {
  const checker = process.platform === "win32" ? "where" : "command";
  const args = process.platform === "win32" ? [command] : ["-v", command];
  return run(checker, args).status === 0;
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

function skillsRoot() {
  const candidates = [
    path.resolve(__dirname, "..", "skills"),
    path.resolve(__dirname, "..", "..", "..", "skills"),
  ];
  for (const candidate of candidates) {
    if (SKILLS.every((name) => fs.existsSync(path.join(candidate, name, "SKILL.md")))) {
      return candidate;
    }
  }
  throw new Error("Could not find bundled atk-jira and atk-cfl skills.");
}

function targetDests(target) {
  const home = os.homedir();
  const codexHome = process.env.CODEX_HOME || path.join(home, ".codex");
  const claudeHome = process.env.CLAUDE_HOME || path.join(home, ".claude");

  if (target === "codex") return [path.join(codexHome, "skills")];
  if (target === "claude") return [path.join(claudeHome, "skills")];
  if (target === "both") return [path.join(codexHome, "skills"), path.join(claudeHome, "skills")];
  throw new Error(`Unknown target: ${target}`);
}

function installSkills(destinations) {
  const root = skillsRoot();
  for (const dest of destinations) {
    for (const name of SKILLS) {
      copyDir(path.join(root, name), path.join(dest, name));
    }
    console.log(`Installed ${SKILLS.join(", ")} skills to ${dest}`);
  }
}

function installCLI() {
  const missing = BINARIES.filter((binary) => !hasCommand(binary));
  if (missing.length === 0) {
    console.log("atk-jira and atk-cfl are already available on PATH.");
    return 0;
  }

  if (hasCommand("brew")) {
    console.log("Installing CLI binaries with Homebrew...");
    const tap = run("brew", ["tap", "wohsj110/tap"], { stdio: "inherit" });
    if (tap.status !== 0) return tap.status || 1;
    for (const cask of ["atk-jira", "atk-cfl"]) {
      const result = run("brew", ["install", "--cask", cask], { stdio: "inherit" });
      if (result.status !== 0) return result.status || 1;
    }
    return 0;
  }

  console.error(`Could not find atk-jira/atk-cfl on PATH and Homebrew is unavailable.

Install the CLI manually from:
  ${REPO}/releases/latest

Then verify:
  atk-jira --help
  atk-cfl --help`);
  return 1;
}

function doctor() {
  let failed = false;
  for (const binary of BINARIES) {
    if (hasCommand(binary)) {
      console.log(`OK: ${binary} is on PATH`);
    } else {
      console.log(`MISSING: ${binary} is not on PATH`);
      failed = true;
    }
  }

  for (const dest of targetDests("both")) {
    for (const skill of SKILLS) {
      const file = path.join(dest, skill, "SKILL.md");
      if (fs.existsSync(file)) {
        console.log(`OK: ${file}`);
      } else {
        console.log(`MISSING: ${file}`);
        failed = true;
      }
    }
  }
  return failed ? 1 : 0;
}

function parseInstallArgs(args) {
  let target = "both";
  let dest = "";
  let shouldInstallCLI = false;

  for (let i = 0; i < args.length; i += 1) {
    const arg = args[i];
    if (arg === "--target") {
      target = args[i + 1];
      i += 1;
    } else if (arg === "--dest") {
      dest = args[i + 1];
      i += 1;
    } else if (arg === "--install-cli") {
      shouldInstallCLI = true;
    } else {
      throw new Error(`Unknown flag: ${arg}`);
    }
  }

  const destinations = dest ? [path.resolve(dest)] : targetDests(target);
  return { destinations, shouldInstallCLI };
}

function main(argv) {
  const command = argv[2];
  if (!command || command === "--help" || command === "-h") {
    usage();
    return 0;
  }

  if (command === "install") {
    const { destinations, shouldInstallCLI } = parseInstallArgs(argv.slice(3));
    installSkills(destinations);
    if (shouldInstallCLI) return installCLI();
    return 0;
  }

  if (command === "add") {
    const name = argv[3];
    if (name && name !== "atlassian-agent" && name !== "atk") {
      throw new Error(`Unknown skill bundle: ${name}`);
    }
    const { destinations, shouldInstallCLI } = parseInstallArgs(argv.slice(name ? 4 : 3));
    installSkills(destinations);
    if (shouldInstallCLI) return installCLI();
    return 0;
  }

  if (command === "install-cli") return installCLI();
  if (command === "doctor") return doctor();

  console.error(`Unknown command: ${command}`);
  usage();
  return 2;
}

try {
  process.exitCode = main(process.argv);
} catch (error) {
  console.error(error.message);
  process.exitCode = 1;
}
