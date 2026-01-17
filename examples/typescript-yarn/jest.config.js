module.exports = {
  preset: "ts-jest",
  reporters: ["default", "jest-junit"],
  testEnvironment: "node",
  testEnvironmentOptions: {
    customExportConditions: ["node", "node-addons"],
  },
};
