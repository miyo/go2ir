module swap_sim();

  reg clk = 1'b0;
  reg reset = 1'b1;
  reg[31:0] counter = 32'h0;
  reg kick = 1'b0;
  wire busy;

  initial begin
    $dumpfile("a.vcd");
    $dumpvars();
  end

  always begin
  #10
  clk <= 1'b0;
  #10
  clk <= 1'b1;
  end

  always @(posedge clk) begin
    counter <= counter + 1;
  end

  always @(posedge clk) begin
    if(counter == 5) reset <= 1'b0;
    if(counter == 7) kick <= 1'b1; else kick <= 1'b0;
    if(counter > 9 && busy == 1'b0) $finish;
  end

  wire[31:0] a = 32'hdeadbeaf, b = 32'habadcafe;
  wire[31:0] result_a, result_b;

  swap U(
    .clk(clk),
    .reset(reset),
    .swap_a(a),
    .swap_b(b),
    .swap_return_0(result_a),
    .swap_return_1(result_b),
    .swap_busy(busy),
    .swap_req(kick)
  );

endmodule
