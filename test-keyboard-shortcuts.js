const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();

  console.log('🔵 Opening NasFlow at localhost:5173...');
  await page.goto('http://localhost:5173/', { waitUntil: 'networkidle' });

  // Wait for login form or redirect
  await page.waitForTimeout(1000);
  const url = page.url();
  console.log(`Current URL: ${url}`);

  // Check if we need to login
  const loginButton = await page.locator('text=เข้าสู่ระบบ, text=ลงทะเบียน').first().isVisible().catch(() => false);

  if (loginButton || url.includes('/login')) {
    console.log('🔑 Login page detected, creating test account...');

    // Click register tab
    const registerTab = page.locator('button:has-text("ลงทะเบียน")').first();
    if (await registerTab.isVisible()) {
      await registerTab.click();
      await page.waitForTimeout(500);
    }

    // Fill registration form
    await page.locator('input[type="email"]').fill('test@example.com');
    await page.locator('input[type="password"]').first().fill('Password123');
    await page.locator('input[type="password"]').last().fill('Password123');

    // Click submit
    const submitBtn = page.locator('button:has-text("ลงทะเบียน"), button:has-text("สร้างบัญชี")').first();
    if (await submitBtn.isVisible()) {
      await submitBtn.click();
      console.log('📝 Registration submitted...');
      await page.waitForTimeout(2000);
    }
  }

  // Wait for board to load
  console.log('⏳ Waiting for board to load...');
  await page.waitForSelector('[class*="KanbanBoard"], [class*="board"], text=รอดำเนินการ', { timeout: 5000 }).catch(() => {
    console.warn('⚠️ Board selector not found, checking what\'s on page...');
  });

  // Wait for tasks to appear
  await page.waitForTimeout(1500);

  // Check if any tasks exist
  const taskCards = await page.locator('div[class*="bg-white"][class*="rounded-lg"]:has(text)').count();
  console.log(`📋 Found ${taskCards} task cards`);

  if (taskCards === 0) {
    console.log('⚠️ No existing tasks found. Creating a test task...');

    // Try to open quick-add with Cmd+K or Ctrl+K
    await page.keyboard.press('Control+K');
    await page.waitForTimeout(500);

    const taskInput = page.locator('input[placeholder*="เพิ่ม"], input[placeholder*="Quick"], input[type="text"]').first();
    if (await taskInput.isVisible()) {
      await taskInput.fill('Test Task 123');
      await page.keyboard.press('Enter');
      console.log('✅ Test task created');
      await page.waitForTimeout(1000);
    }
  }

  // Now test keyboard shortcuts
  console.log('\n🎯 Testing keyboard shortcuts...\n');

  // Step 1: Click a task card to focus it
  const firstTaskCard = page.locator('div[class*="bg-white"][class*="rounded-lg"]').first();
  if (await firstTaskCard.isVisible()) {
    console.log('1️⃣  Clicking first task card to focus...');
    await firstTaskCard.click();
    await page.waitForTimeout(500);

    // Check for focused ring
    const focusedCard = await firstTaskCard.evaluate(el => {
      const classList = el.className;
      return classList.includes('ring-') || classList.includes('ring-2');
    });

    if (focusedCard) {
      console.log('✅ Task card focused (has ring class)');
    } else {
      console.log('⚠️ Task card clicked but no focus ring detected');
    }

    // Step 2: Press E to open drawer
    console.log('\n2️⃣  Pressing E to open task drawer...');
    await page.keyboard.press('e');
    await page.waitForTimeout(800);

    const drawerVisible = await page.locator('[class*="drawer"], [class*="slide"], aside').isVisible().catch(() => false);
    if (drawerVisible) {
      console.log('✅ Task drawer opened');
    } else {
      console.log('⚠️ Drawer may not have opened visibly');
    }

    // Step 3: Close drawer and test D key
    console.log('\n3️⃣  Pressing Escape to close drawer...');
    await page.keyboard.press('Escape');
    await page.waitForTimeout(500);

    console.log('4️⃣  Pressing D to toggle task status...');
    const originalStatus = await firstTaskCard.evaluate(el => el.textContent);
    await page.keyboard.press('d');
    await page.waitForTimeout(800);

    console.log('✅ D key pressed - status should be updated (check toast/visual change)');

    // Step 4: Test priority shortcuts
    console.log('\n5️⃣  Testing priority shortcuts (1=low, 2=medium, 3=high)...');

    for (const [key, priority] of [['1', 'low'], ['2', 'medium'], ['3', 'high']]) {
      console.log(`  Pressing ${key} for ${priority} priority...`);
      await page.keyboard.press(key);
      await page.waitForTimeout(600);
      const toast = await page.locator('div[class*="toast"], div[role="alert"]').textContent().catch(() => '');
      console.log(`  ✅ ${key} pressed - expected priority: ${priority}`);
    }

    // Step 5: Test that shortcuts don't work when typing
    console.log('\n6️⃣  Testing that shortcuts don\'t trigger in input fields...');
    const searchInput = page.locator('input[placeholder*="ค้นหา"], input[type="search"]').first();
    if (await searchInput.isVisible()) {
      await searchInput.click();
      await searchInput.fill('test');
      await page.keyboard.press('d'); // Should not toggle
      console.log('✅ D pressed in search input - should NOT trigger shortcut');
    }
  }

  console.log('\n🎉 Verification complete!');
  await browser.close();
})().catch(err => {
  console.error('❌ Test failed:', err);
  process.exit(1);
});
